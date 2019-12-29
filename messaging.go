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

package gorge

import (
	"reflect"
	"strings"
	"sync"
	"time"
)

// Messaging thing
type Messaging struct {
	parent *Messaging
	Groups map[reflect.Type]*HandlerGroup
	links  map[*Messaging]struct{}
}

// Link will add child messagers
// TODO: maybe, check parent persistences and trigger here
func (m *Messaging) Link(sub *Messaging) {
	if m.links == nil {
		m.links = map[*Messaging]struct{}{}
	}

	sub.parent = m
	m.links[sub] = struct{}{}
}

//Unlink remove a sub messaging system
func (m *Messaging) Unlink(sub *Messaging) {
	if m.links == nil {
		return
	}
	delete(m.links, sub)
}

// ReflAuto registers func handelrs by checking Methods on a struct and the
// method signature
// If method has Handler prefix it will register as an handler
// If method has Watch prefix it will register as a watcher
func (m *Messaging) ReflAuto(s interface{}) {
	f := reflect.ValueOf(s)
	sn := f.Elem().Type().Name()
	for i := 0; i < f.NumMethod(); i++ {
		mn := f.Type().Method(i).Name
		if strings.HasPrefix(mn, "handle") {
			m.Handle(f.Method(i).Interface()).Describe(sn + "." + mn)
			continue
		}
	}
}

// Handle registers an handle that triggers upon registration if there is a last
// message, it will return the registered handler
// the fn interface must be a function with 1 parameter and no returns
// (i.e: func(type))
func (m *Messaging) Handle(fn interface{}) *Handler {
	k := fnTyp(fn)
	g := m.groupOrNew(k)

	h := g.AddFn(fn)

	if val := m.lastValue(k); val != nil {
		args := []reflect.Value{reflect.ValueOf(val)}
		reflect.ValueOf(fn).Call(args)
	}
	return h
}

// Query will callback if the signature type has a message
// the fn interface must be a function with 1 parameter and no returns
// (i.e: func(type))
func (m *Messaging) Query(fn interface{}) {
	val := m.lastValue(fnTyp(fn))
	if val == nil {
		return
	}
	args := []reflect.Value{reflect.ValueOf(val)}
	reflect.ValueOf(fn).Call(args)
}

// Trigger will retrieve handlers for the specific signature type and call each
// one
//
// XXX: Might be bad to hold a reference to last value while triggering since
// it can be a huge value, consider creating a new kind of Trigger method for this
// purpose like `Persist(v interface)`
func (m *Messaging) Trigger(v interface{}) {
	k := reflect.TypeOf(v)
	// Trigger local
	if g := m.group(k); g != nil {
		g.Call(v)
	}
	m.trigger(k, v)
}

// Persist works like trigger but persists the message for new Watchers
func (m *Messaging) Persist(v interface{}) *HandlerGroup {
	k := reflect.TypeOf(v)

	g := m.groupOrNew(k)
	g.Last = v
	g.Call(v)

	m.trigger(k, v)

	return g
}
func (m *Messaging) trigger(k reflect.Type, v interface{}) {
	m.triggerParent(k, v)
	m.triggerLinks(k, v)
}
func (m *Messaging) triggerParent(k reflect.Type, v interface{}) {
	if m.parent == nil {
		return
	}
	if g := m.parent.group(k); g != nil {
		g.Call(v)
	}
	// Go up
	m.parent.triggerParent(k, v)
}
func (m *Messaging) triggerLinks(k reflect.Type, v interface{}) {
	if m.links == nil {
		return
	}
	for l := range m.links {
		if g := l.group(k); g != nil {
			g.Call(v)
		}
		l.triggerLinks(k, v)
	}
}

func (m *Messaging) lastValue(k reflect.Type) interface{} {
	cur := m
	for cur != nil {
		if g := cur.group(k); g != nil && g.Last != nil {
			return g.Last
		}
		cur = cur.parent
	}
	return nil
}

// entryFor retrieves the entry for the type, if create is true it will create it
func (m *Messaging) group(k reflect.Type) *HandlerGroup {
	if m.Groups == nil {
		return nil
	}
	return m.Groups[k]
}
func (m *Messaging) groupOrNew(k reflect.Type) *HandlerGroup {
	if m.Groups == nil {
		m.Groups = map[reflect.Type]*HandlerGroup{}
	}
	g, ok := m.Groups[k]
	if !ok {
		g = &HandlerGroup{
			Type:     k,
			Handlers: []*Handler{},
		}
		m.Groups[k] = g
	}
	return g
}

// Handler holds the func callback and some extra fields for profiling
type Handler struct {
	Desc      string
	CallStart time.Time
	CallEnd   time.Time
	Once      bool
	Fn        interface{}
}

// Describe the handler for debug purposes
func (h *Handler) Describe(v string) *Handler {
	h.Desc = v
	return h
}

// HandlerGroup holds the list of handlers for a specific function signature
type HandlerGroup struct {
	sync.Mutex
	Type      reflect.Type
	Last      interface{}
	CallStart time.Time
	CallEnd   time.Time
	Handlers  []*Handler
}

// AddFn creates an handler with the fn and adds
func (g *HandlerGroup) AddFn(fn interface{}) *Handler {
	return g.Add(&Handler{Fn: fn})
}

// Add the handler to entry
func (g *HandlerGroup) Add(h *Handler) *Handler {
	g.Lock()
	g.Handlers = append(g.Handlers, h)
	g.Unlock()

	return h
}

// Remove handler
func (g *HandlerGroup) Remove(h *Handler) {
	g.Lock()
	defer g.Unlock()
	for i, lh := range g.Handlers {
		if lh == h {
			g.Handlers = append(g.Handlers[:i], g.Handlers[i+1:]...)
			break
		}
	}
}

// Range the handlers
func (g *HandlerGroup) Range(fn func(h *Handler) bool) {
	g.Lock()
	defer g.Unlock()

	for _, h := range g.Handlers {
		if !fn(h) {
			break
		}
	}
}

// Call calls the funcs on the handlers
func (g *HandlerGroup) Call(v interface{}) {
	args := []reflect.Value{reflect.ValueOf(v)}
	g.Lock()
	defer g.Unlock()
	if g.Handlers == nil {
		return
	}
	// Removable Once handlers
	var toRemove []int
	g.CallStart = time.Now()
	for i, h := range g.Handlers {
		h.CallStart = time.Now()
		reflect.ValueOf(h.Fn).Call(args)
		h.CallEnd = time.Now()
		if h.Once {
			toRemove = append(toRemove, i)
		}
	}
	g.CallEnd = time.Now()
	for _, hi := range toRemove {
		g.Handlers = append(g.Handlers[:hi], g.Handlers[hi+1:]...)
	}
}

// fnTyp retrieves the dominant type for the specific signature func(type)
func fnTyp(fn interface{}) reflect.Type {
	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func && typ.NumIn() != 1 {
		panic("wrong type, should be a func with 1 param")
	}
	return typ.In(0)
}
