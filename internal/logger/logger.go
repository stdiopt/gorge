// Package logger overrides to Override the writer on log.SetOutput
package logger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type prettylogStyle struct {
	Counter  Style
	Message  Style
	Prefix   Style
	Time     Style
	Duration Style
	File     Style
}

var (
	// GlobalStyle style for each log area.
	GlobalStyle = prettylogStyle{
		Counter:  Style{Prefix: "\033[37m", Suffix: "\033[0m", IncrementPad: true},
		Message:  Style{Prefix: "\033[37m", Suffix: "\033[0m"},
		Prefix:   Style{IncrementPad: true},
		Time:     Style{Prefix: "\033[34m", Suffix: "\033[0m"},
		Duration: Style{Prefix: "\033[90m", Suffix: "\033[0m"},
		File:     Style{Prefix: "\033[30m", Suffix: "\033[0m"},
	}

	// Color per Prefix.
	prefixStyle  = map[string]int{}
	prefixColors = []string{
		"\033[31m", "\033[32m", "\033[33m", "\033[35m", "\033[36m",
		"\033[01;31m", "\033[01;32m", "\033[01;33m", "\033[01;34m",
		"\033[01;35m", "\033[01;36m",
	}
)

// Writer writer struct.
type Writer struct {
	prefix   string
	lastTime time.Time
	counter  int64
	writers  []io.Writer
}

// NewWriter creates a new log writer to be used in log.SetOutput.
func NewWriter(prefix string, writers ...io.Writer) *Writer {
	wri := writers
	if len(wri) == 0 { // defaults to stderr if none
		wri = []io.Writer{os.Stderr}
	}
	return &Writer{prefix, time.Now(), 0, wri}
}

// Write io.Write implementation that parses the output.
func (w *Writer) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	originalLen := len(b)

	written := 0
	for {
		n := bytes.IndexByte(b, '\n')
		if n == -1 {
			if err := w.writeLine(b); err != nil {
				return written, err
			}
			break
		}
		if err := w.writeLine(b[:n]); err != nil {
			return written, err
		}
		written += len(b[:n])
		b = b[n+1:]
		if len(b) == 0 {
			break
		}
	}
	return originalLen, nil
}

func (w *Writer) writeLine(b []byte) error {
	// Writes the line

	msg := string(b)
	var ptr uintptr
	var file string
	var line int
	var ok bool
	var tname string
	for call := 3; call < 5; call++ {
		ptr, file, line, ok = runtime.Caller(call)
		if !ok {
			continue
		}
		tname = runtime.FuncForPC(ptr).Name()
		if !strings.HasPrefix(tname, "log") {
			break
		}
	}

	pkgMethod := tname[strings.LastIndex(tname, "/")+1:]
	pkg := pkgMethod[:strings.Index(pkgMethod, ".")]
	fname := file[strings.LastIndex(file, "/")+1:]

	timeDiff := time.Since(w.lastTime)

	duration := durationStr(timeDiff)

	prefixStr := pkg
	if w.prefix != "" {
		prefixStr = w.prefix
	}
	// Colored prefix, it will match a string in the prefix map
	// and fetch correspondent color in the color list
	prefixStyleID, ok := prefixStyle[prefixStr]
	if !ok {
		prefixStyleID = len(prefixStyle)
		prefixStyle[prefixStr] = prefixStyleID
	}
	prefixColor := prefixColors[prefixStyleID%len(prefixColors)]

	str := fmt.Sprintf("[%s:%s %s]: %s %s %s\n",
		GlobalStyle.Counter.Get(w.counter),
		GlobalStyle.Time.Get(time.Now().Format("2006-01-02 15:04:05.000")),
		GlobalStyle.Prefix.GetCustom(prefixColor, "\033[0m", prefixStr),
		GlobalStyle.Message.Get(msg),

		GlobalStyle.Duration.Get(duration),
		GlobalStyle.File.Get(fmt.Sprintf("%s:%d", fname, line)),
	)

	w.lastTime = time.Now()
	w.counter++

	for _, ww := range w.writers {
		_, err := ww.Write([]byte(str))
		if err != nil {
			return err
		}
	}
	return nil
}

// New creates a new log.Logger with a prefix.
func New(prefix string, writers ...io.Writer) *log.Logger {
	return log.New(NewWriter(prefix, writers...), "", 0)
}

// Dummy a log.Logger with io.Discard writer.
func Dummy() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}

// Global sets the global log with a prettylog writer.
func Global() {
	log.SetFlags(0)
	log.SetOutput(NewWriter(""))
}

func durationStr(dur time.Duration) string {
	fdurationSuf := "ms"
	fduration := float64(dur.Nanoseconds()) / 1000000.0
	if fduration > 100 {
		fduration /= 1000
		fdurationSuf = "s"
	}

	return fmt.Sprintf("+%.2f/%s", fduration, fdurationSuf)
}
