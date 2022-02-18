// Package setlist slice without duplication.
package setlist

type SetList[T comparable] struct {
	items []T
	uniq  map[T]struct{}
}

func (l *SetList[T]) Items() []T {
	return l.items
}

func (l *SetList[T]) Add(item T) bool {
	if l.uniq == nil {
		l.uniq = make(map[T]struct{})
	}
	if _, ok := l.uniq[item]; ok {
		return false
	}

	l.items = append(l.items, item)
	l.uniq[item] = struct{}{}
	return true
}

func (l *SetList[T]) Remove(item T) bool {
	if _, ok := l.uniq[item]; !ok {
		return false
	}

	delete(l.uniq, item)

	for i, v := range l.items {
		if v == item {
			t := l.items
			l.items = append(l.items[:i], l.items[i+1:]...)
			var z T
			// zero value to unreference any pointers.
			t[len(t)-1] = z
			return true
		}
	}
	return false
}

// First will panic if the list is empty.
func (l *SetList[T]) Front() T {
	return l.items[0]
}

func (l *SetList[T]) Len() int {
	return len(l.items)
}
