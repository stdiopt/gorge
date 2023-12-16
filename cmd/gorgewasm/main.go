package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed tmpl
var tmplFS embed.FS
var indexTMPL = template.Must(template.ParseFS(tmplFS, "tmpl/index.tmpl"))

func main() {
	log.SetFlags(0)
	err := run(os.Args[1:])
	if err != nil {
		log.Fatal("err", err)
	}
}

func usage(fl *flag.FlagSet) {
	usage := []string{
		"Usage:",
		"",
		"\tbuild <pkg>",
		"\tserve <pkg>",
		"",
	}
	fmt.Fprintln(os.Stderr, strings.Join(usage, "\n"))

	if fl != nil {
		fl.Usage()
	}
}

func run(args []string) error {
	if len(args) < 1 {
		usage(nil)
		return errors.New("missing arguments")
	}

	switch args[0] {
	case "build":
		output := "wasm.html"

		fl := flag.NewFlagSet(args[0], flag.ExitOnError)
		fl.StringVar(&output, "o", "wasm.html", "output file")
		if err := fl.Parse(args[1:]); err != nil {
			return err
		}

		pkg := fl.Arg(0)
		if pkg == "" {
			fl.Usage()
			return fmt.Errorf("build: missing arg <pkg>")
		}

		log.Printf("Writing.. %s", output)
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := buildHTML(f, pkg); err != nil {
			return err
		}

	case "serve":
		laddr := ":8080"
		fl := flag.NewFlagSet(args[0], flag.ExitOnError)
		fl.StringVar(&laddr, "l", ":8080", "address to listen http")
		if err := fl.Parse(args[1:]); err != nil {
			return err
		}

		pkg := fl.Arg(0)
		if pkg == "" {
			fl.Usage()
			return fmt.Errorf("serve: missing arg <pkg>")
		}
		log.Println("Listening at:", laddr)

		// Helper
		if strings.HasPrefix(laddr, ":") {
			port := laddr[1:]
			addrW := bytes.NewBuffer(nil)
			fmt.Fprintf(addrW, "    http://localhost:%s\n", port)
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				log.Fatal("err:", err)
			}
			for _, a := range addrs {
				astr := a.String()
				if strings.HasPrefix(astr, "192.168") ||
					strings.HasPrefix(astr, "10") {
					a := strings.Split(astr, "/")[0]
					fmt.Fprintf(addrW, "    http://%s:%s\n", a, port)
				}
			}
			log.Println(addrW.String())
		}

		return http.ListenAndServe(laddr, wasmServeHTTP(pkg))
	default:
		usage(nil)
	}
	return nil
}

func wasmServeHTTP(p string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
			return
		}
		if err := buildHTML(w, p); err != nil {
			http.Error(w, "error rendering html", http.StatusInternalServerError)
			return
		}
	}
}

func buildHTML(w io.Writer, p string) error {
	wasmBuf := &bytes.Buffer{}

	err := buildWasm(wasmBuf, p)
	if err != nil {
		return fmt.Errorf("wasm build error: %w", err)
	}
	err = indexTMPL.Execute(w, map[string]any{
		"wasmexec": template.JS(wasmExec),
		"wasmcode": base64.StdEncoding.EncodeToString(wasmBuf.Bytes()),
	})
	if err != nil {
		return fmt.Errorf("html build error: %w", err)
	}
	return nil
}

func buildWasm(w io.Writer, pkg string) error {
	log.Printf("building %v...", pkg)
	tf, err := os.CreateTemp(os.TempDir(), "gorgewasm.")
	if err != nil {
		return err
	}
	if err := tf.Close(); err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(tf.Name()) // nolint: errcheck
		if err != nil {
			log.Println("error removing temp files:", err)
		}
	}()

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
	_, err = io.Copy(w, f)
	return err
}

var wasmExec = func() []byte {
	goroot := build.Default.GOROOT
	wasmExecName := filepath.Join(goroot, "misc/wasm/wasm_exec.js")
	// Read wasm_exec from system dist
	data, err := os.ReadFile(wasmExecName)
	if err != nil {
		panic(err)
	}
	return data
	// return bytes.ReplaceAll(data, []byte("console.log"), []byte("myLog"))
}()
