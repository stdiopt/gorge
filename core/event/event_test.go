package event_test

import (
	"testing"

	"github.com/stdiopt/gorge/core/event"
)

type B struct {
	v int
}

func (b *B) HandleEvent(v event.Event) { b.v++ }

func TestRemoveHandler(t *testing.T) {
	b := &event.Bus{}

	h := &B{}
	b.Handle(h)

	b.Trigger(1)

	b.RemoveHandler(h)

	b.Trigger(1)

	if want := 1; h.v != want {
		t.Errorf("\nwant: %v\n got: %v\n", want, h.v)
	}
}
