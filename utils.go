package gorge

// StringHash builds an hash from a string
func StringHash(str ...string) uint {
	seed := uint(0)
	for _, s := range str {
		for _, c := range s {
			seed = 31*seed + uint(c)
		}
	}
	return seed
}

// EachEntity solves entities containers and walk on every entity.
func EachEntity(e Entity, fn func(e Entity)) {
	fn(e)
	if v, ok := e.(EntityContainer); ok {
		for _, e := range v.GetEntities() {
			EachEntity(e, fn)
		}
	}
}

// EachParent iterates parents
func EachParent(e Entity, fn func(e Entity) bool) {
	for e != nil {
		if !fn(e) {
			return
		}
		p, ok := e.(ParentGetter)
		if !ok {
			break
		}
		e = p.Parent()
	}
}

// HasParent verifies if the parent exists in e hierarchy
func HasParent(e Entity, parent Entity) bool {
	hasParent := false
	EachParent(e, func(e Entity) bool {
		if e == parent {
			hasParent = true
			return false
		}
		return true
	})
	return hasParent
}
