// Package setlist is an append orderer set.
// it has better performance than slice while removing elements
// and avoid duplications on insert
package setlist

import "container/list"

// SetList is an append ordered Set.
type SetList struct {
	list  list.List
	index map[interface{}]*list.Element
}

// Add to list, returns true if element was added, false if already existed
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

// Get returns the value based on k, returns true if exists, false otherwise.
func (l *SetList) Get(k interface{}) (interface{}, bool) {
	if l.index == nil {
		return nil, false
	}
	r, ok := l.index[k]
	return r, ok
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

// Front returns the first element of the list
func (l *SetList) Front() *list.Element {
	return l.list.Front()
}

// Len returns the number of items
func (l *SetList) Len() int {
	return l.list.Len()
}
