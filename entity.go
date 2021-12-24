package gorge

import (
	"bytes"
	"fmt"
)

// Entity can be anything but shouldn't be a func
type Entity interface{}

// EntityContainer is an interface used while adding entity to solve EventAddEntity
// one container with 2 entities will trigger 2 EventAddEntity.
type EntityContainer interface {
	GetEntities() []Entity
}

/*
// EntityFunc a type of func with an entity as argument.
type EntityFunc func(e Entity)

// ApplyTo apply the func to an entity, used to satisfy interface.
func (fn EntityFunc) ApplyTo(e Entity) { fn(e) }

// FuncGroup returns an EntityFunc that calls a slice of entity funcs.
func FuncGroup(opts ...EntityFunc) EntityFunc {
	return func(e Entity) {
		for _, fn := range opts {
			fn(e)
		}
	}
}


// ApplyTo executes several entity funcs on an entity.
func ApplyTo(e Entity, opts ...EntityFunc) {
	FuncGroup(opts...).ApplyTo(e)
}
*/

// Container provides a way to multiplex entities by implementing EntityContainer interface.
type Container []Entity

// Enforce interface
var _ EntityContainer = (*Container)(nil)

// Add one ore more entities to Container
// it will loop existing elements to compare and avoid duplicates.
func (c *Container) Add(ents ...Entity) {
	for _, ent := range ents {
		if c.indexOf(ent) != -1 {
			continue
		}
		*c = append(*c, ent)
	}
}

// Remove entities from container it does not remove from the world.
func (c *Container) Remove(ents ...Entity) {
	for _, ent := range ents {
		n := c.indexOf(ent)
		if n == -1 {
			continue
		}
		t := (*c)
		*c = append((*c)[:n], (*c)[n+1:]...)
		// deref last entity to be able to GC.
		t[len(t)-1] = nil
		break
	}
}

func (c *Container) indexOf(ent Entity) int {
	for i, e := range *c {
		if e == ent {
			return i
		}
	}
	return -1
}

// GetEntities implements the entitycontainer
func (c Container) GetEntities() []Entity { return c }

func (c Container) String() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprint(buf, "[")
	for i, e := range c {
		if i != 0 {
			fmt.Fprint(buf, ", ")
		}
		fmt.Fprintf(buf, "%T", e)
	}
	fmt.Fprint(buf, "]")
	return buf.String()
}
