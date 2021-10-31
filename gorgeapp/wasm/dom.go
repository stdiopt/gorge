//go:build (js && ignore) || wasm

package wasm

import "syscall/js"

// Browser dom stuff
var (
	Document = js.Global().Get("document")
	Body     = Document.Get("body")
)

// Util
type (
	// Class will tell El() to add these into the classList
	Class []string
	// Attr will tell El() to set attributes
	Attr map[string]string
	// Text will tell El() to create a text node
	Text string
)

// El create a Dom element
func El(name string, args ...interface{}) js.Value {
	e := Document.Call("createElement", name)
	for _, a := range args {
		switch va := a.(type) {
		case Class:
			for _, v := range va {
				e.Get("classList").Call("add", v)
			}
		case Attr:
			for k, v := range va {
				e.Call("setAttribute", k, v)
			}
		case string:
			e.Set("innerHTML", e.Get("innerHTML").String()+va)
		case Text:
			e.Call("appendChild", Document.Call("createTextNode", string(va)))
		case js.Value:
			e.Call("appendChild", va)
		default:
			panic("not sure what todo")
		}
	}
	return e
}
