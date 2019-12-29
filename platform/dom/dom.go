// Copyright 2019 Luis Figueiredo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build js, wasm

package dom

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
