package gorge

import "github.com/stdiopt/gorge/core/event"

type gorger interface {
	Gorge() *Context
}

func TriggerInMain[T any](g gorger, e T) {
	g.Gorge().RunInMain(func() {
		Trigger(g.Gorge(), e)
	})
}

func Trigger[T any](b event.Buser, v T) {
	event.Trigger(b, v)
}

func HandleFunc[T any](b event.Buser, fn func(T)) {
	event.HandleFunc(b, fn)
}

// EventPreUpdate type
type EventPreUpdate float32

// DeltaTime returns the float32 delta time for the event.
func (e EventPreUpdate) DeltaTime() float32 { return float32(e) }

// EventUpdate type
type EventUpdate float32

// DeltaTime returns the float32 delta time for the event.
func (e EventUpdate) DeltaTime() float32 { return float32(e) }

// EventPostUpdate type
type EventPostUpdate float32

// DeltaTime returns the float32 delta time for the event.
func (e EventPostUpdate) DeltaTime() float32 { return float32(e) }

// EventRender happens after pre,update and post update events
type EventRender float32

// EventAddEntity is triggered when entities are added
type EventAddEntity struct {
	Entity
}

// EventRemoveEntity is triggered when entities are destroyed
type EventRemoveEntity struct {
	Entity
}

// EventStart fired when things starts
type EventStart struct{}

// EventAfterStart to attach stuff (wasm request animation frame workaround)
type EventAfterStart struct{}

// EventDestroy is called when system is shutting down
type EventDestroy struct{}

// EventError contains an error
type EventError struct{ Err error }

// EventWarn contains a warning
type EventWarn string

// EventResourceUpdate sends a resource through systems for aditional treatment
// i.e: uploading to gpu
type EventResourceUpdate struct {
	Resource any
}

/*func UpdateHandler(fn func(e EventUpdate)) event.HandlerFunc {
	return func(v event.Event) {
		if e, ok := v.(EventUpdate); ok {
			fn(e)
		}
	}
}*/

/*func HandleUpdate(fn func(e EventUpdate)) {
	HandleFunc(UpdateHandler(fn))
}*/
