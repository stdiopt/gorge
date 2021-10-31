// This will compile wasm package and serve files
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/gohxs/prettylog"
)

var (
	tmpl         *template.Template
	wasmExec     []byte
	flagHTMLFile string
)

func init() {
	var err error
	wasmExecName := filepath.Join(build.Default.GOROOT, "misc/wasm/wasm_exec.js")
	// Read wasm_exec from system dist
	wasmExec, err = ioutil.ReadFile(wasmExecName)
	if err != nil {
		panic(err)
	}
	t := template.New("wasm")
	tmpl, err = t.Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
}

func main() {
	prettylog.Global()
	flag.StringVar(&flagHTMLFile, "ohtml", "", "write output to file.wasm file.html")
	flag.Parse()
	pkg := flag.Arg(0)

	if flagHTMLFile != "" {
		if err := writeHTMLFile(flagHTMLFile, pkg, "main.wasm"); err != nil {
			log.Fatal(err)
		}
		return
	}
	var listener net.Listener
	port := 8080
	var addr string
	for {
		addr = fmt.Sprintf(":%d", port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Println("Err opening", port, err)
			port++
			log.Println("Trying port", port)
			continue
		}
		listener = lis
		break
	}

	log.Println("Listening at:", addr)
	err := http.Serve(listener, chain(
		Server(),
		Logger(prettylog.New("gorgeserve")),
	))
	if err != nil {
		log.Fatal(err)
	}
}

// This html contains a loader that informs the wasm loaded progress and a
// template to be injected in an iframe, the inner iframe will load wasm file
// and start it, this way any it's safe for the wasm code to mess with dom without
// compromising the main loader elements.
var htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
<template id='wasm_script' type="text/raw">
	<body>wasm</body>
	<script>
	{{ .wasmexec }}
	const nativeLog = console.log;
	console.log = (...args) => {
		parent.postMessage({type: "output", msg: args.join(" ")}, '*')
		nativeLog(...args)
	}
	async function loader(wasmFile) {
		try {
			const res = await fetch(wasmFile)
			if (res.status != 200) { throw await res.text(); }
			const reader = res.body.getReader();
			const total = res.headers.get('content-length')
			let bytes = new Uint8Array(total)
			for(let cur=0;;) {
				const {done, value} = await reader.read();
				if (done || !value) { break }
				bytes.set(value, cur)
				cur += value.length
				parent.postMessage({type: "progress", msg: (cur / total)})
			}
			parent.postMessage({type: "progress", msg: "done"})
			document.body.innerHTML = ""
			const go = new Go()
			await go.run((await WebAssembly.instantiate(bytes.buffer, go.importObject)).instance)
			parent.postMessage({type: "done"})
		} catch(err) {
			console.log(err)
			parent.postMessage({type: "error"})
		}
	}
	window.addEventListener('message', evt => loader(evt.data))
	window.addEventListener("keyup", (evt) => evt.key == "Escape" && parent.postMessage({type: "switch"}) )
	</script>
</template>

<script>
const toggler = (...args) => (active) => args.forEach(v => {v.style.display=v==active?'':'none';v.focus()})
window.onload = function(evt) {
	const $loader	= document.querySelector("loader")
	const $iframe	= document.querySelector("iframe")
	const $counter	= document.querySelector("counter")
	const $progress	= document.querySelector("progressvalue")
	const $output	= document.querySelector("output")
	const toggle = toggler($loader,$iframe)
	window.addEventListener("keyup", (evt) => evt.key == "Escape" && toggle($iframe))
	window.addEventListener('message', (evt) => {
		switch (evt.data.type) {
			case "progress":
				if (evt.data.msg == "done"){toggle($iframe);break}
				const percent = (evt.data.msg * 100).toFixed(2) + '%'
				$progress.style.width = percent
				$counter.innerHTML = "Receiving&nbsp;<i>{{.pkg}}</i>&nbsp;<b>"+percent+"</b>"
				break
			case "done": case "error": case "switch": 
				toggle($loader);break
			default:
				$output.innerHTML += evt.data.msg+"\n"
		}
	})
	$iframe.contentDocument.write( document.querySelector("#wasm_script").innerHTML)
	const wasmFile	= (new URLSearchParams(window.location.search)).get("t") || "main.wasm"
	$iframe.contentWindow.postMessage(wasmFile)
}
</script>

<style>
body,body *{box-sizing:border-box;display:flex;border:none;}
body{position:relative;margin:0;padding:0;height:100vh;flex-flow:column;align-items:stretch;}body>*{flex:1;}
loader{position:relative;padding:20px;flex-flow:column;justify-content:center;align-items:center;}
progressvalue{background:#00add8;box-shadow:0px 0px 5px 0px cyan;width:0%;}
progressbar{margin:10px;height:3px;width:200px;justify-content:center;background:rgba(0,200,255,0.1);}
output{flex-flow:column;white-space:pre-wrap;font-family:monospace;max-width:80em;overflow-wrap:anywhere;color:#aa4400;}
</style>
</head>
<body>
<loader><counter></counter><progressbar><progressvalue></progressvalue></progressbar><output></output></loader>
<iframe style="display:none"/>
</body>
`

// Server http handler to serve wasm
func Server() http.HandlerFunc {
	pkg := flag.Arg(0)
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		switch path {
		case "/":
			if err := buildHTML(w, pkg, "main.wasm"); err != nil {
				writeStatus(w, http.StatusInternalServerError, err)
			}
		case "/main.wasm":
			if err := buildWasm(w, pkg); err != nil {
				writeStatus(w, http.StatusInternalServerError, err)
			}
		default:
			path = path[1:]
			http.ServeFile(w, r, path)
		}
	}
}

func writeHTMLFile(fname, pkg, wasmFile string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	return buildHTML(f, pkg, wasmFile)
}

func buildHTML(w io.Writer, pkg, wasmFile string) error {
	if w, ok := w.(http.ResponseWriter); ok {
		w.Header().Set("Content-type", "text/html")
	}
	topts := map[string]interface{}{
		"pkg":      pkg,
		"wasmexec": string(wasmExec),
		"wasmfile": wasmFile,
	}
	return tmpl.ExecuteTemplate(w, "wasm", topts)
}

func buildWasm(w http.ResponseWriter, pkg string) error {
	log.Printf("building %v...", pkg)
	tf, err := ioutil.TempFile(os.TempDir(), "gorgewasm.")
	if err != nil {
		return err
	}
	if err := tf.Close(); err != nil {
		return err
	}
	defer os.Remove(tf.Name()) // nolint: errcheck

	versionCmd := exec.Command("go", "version")
	versionCmd.Stdout = log.Writer()
	if err := versionCmd.Run(); err != nil {
		return err
	}

	// BUILDCOMMAND
	errBuf := new(bytes.Buffer)
	cmd := exec.Command("go", "build", "-o", tf.Name(), pkg)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, errBuf)
	if err := cmd.Run(); err != nil {
		return errors.New(errBuf.String())
	}

	f, err := os.Open(tf.Name())
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck
	oi, err := f.Stat()
	if err != nil {
		return err
	}
	w.Header().Set("Content-type", "application/wasm")
	w.Header().Set("Content-length", fmt.Sprint(oi.Size()))
	_, err = io.Copy(w, f)
	return err
}

type logHelper struct {
	http.ResponseWriter
	statusCode int
}

func (l *logHelper) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

// Logger middleware
func Logger(log *log.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := &logHelper{w, 200}
			if next != nil {
				next.ServeHTTP(l, r)
			}
			raddr := r.RemoteAddr
			log.Printf("(%s) %s %s - [%d %s]", raddr, r.Method, r.URL.Path, l.statusCode, http.StatusText(l.statusCode))
		})
	}
}

// MiddlewareFunc type for middleware chain
type MiddlewareFunc func(http.Handler) http.Handler

func chain(next http.Handler, mws ...MiddlewareFunc) http.Handler {
	if len(mws) == 0 {
		return next
	}
	return mws[0](chain(next, mws[1:]...))
}

func writeStatus(w http.ResponseWriter, code int, extras ...interface{}) {
	w.WriteHeader(code)
	extra := fmt.Sprint(extras...)
	fmt.Fprint(w, http.StatusText(code), "\n", extra)
}
