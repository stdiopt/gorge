package gorge

import "container/list"

//SetList is an append ordered Set.
// it has better performance than slice while removing elements
// and avoid duplications on insert
type SetList struct {
	list  list.List
	index map[interface{}]*list.Element
}

//Add to list, returns true if element was added, false if already existed
func (l *SetList) Add(v interface{}) bool {
	if l.index == nil {
		l.index = map[interface{}]*list.Element{}
	}
	if _, ok := l.index[v]; ok {
		return false
	}
	e := l.list.PushBack(v)
	l.index[v] = e
	return true
}

// Remove from list, returns true if removed, false if not found
func (l *SetList) Remove(v interface{}) bool {
	if l.index == nil {
		return false
	}
	e, ok := l.index[v]
	if !ok {
		return false
	}
	l.list.Remove(e)
	delete(l.index, v)
	return true
}

// Range the list
func (l *SetList) Range(fn func(v interface{}) bool) {
	var next *list.Element
	for e := l.list.Front(); e != nil; e = next {
		next = e.Next()
		if !fn(e.Value) {
			return
		}
	}
}

// Len returns the number of items
func (l *SetList) Len() int {
	return l.list.Len()
}
