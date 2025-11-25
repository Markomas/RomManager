package util

type OrderedSet[T comparable] struct {
	items []T
	index map[T]int // maps item → index in slice
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		items: []T{},
		index: make(map[T]int),
	}
}

// Add adds an item if it does not exist
func (s *OrderedSet[T]) Add(item T) {
	if _, exists := s.index[item]; exists {
		return
	}
	s.items = append(s.items, item)
	s.index[item] = len(s.items) - 1
}

// Exists checks if the item exists
func (s *OrderedSet[T]) Exists(item T) bool {
	_, ok := s.index[item]
	return ok
}

// Last returns the last inserted item
func (s *OrderedSet[T]) Last() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Remove removes an item (O(1))
func (s *OrderedSet[T]) Remove(item T) {
	i, exists := s.index[item]
	if !exists {
		return
	}

	lastIndex := len(s.items) - 1
	lastItem := s.items[lastIndex]

	// swap with last to achieve O(1) delete
	s.items[i] = lastItem
	s.index[lastItem] = i

	// shrink slice
	s.items = s.items[:lastIndex]

	// remove from map
	delete(s.index, item)
}

func (s *OrderedSet[T]) MoveToFront(path T) {
	_, exists := s.index[path]
	if !exists {
		return
	}

	s.Remove(path)
	s.Add(path)
}
