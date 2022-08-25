// Package hashset provides an implementation of a set using the built-in map.
package hashset

// New returns an empty hashset.
func New[T comparable](values ...T) Set[T] {
	set := Set[T]{
		values: make(map[T]struct{}),
	}

	set.Add(values...)

	return set
}

// Set implements a hashset, using the hashmap as the underlying storage.
type Set[T comparable] struct {
	values map[T]struct{}
}

// Add adds 'values' to the set.
func (s Set[T]) Add(values ...T) {
	for _, v := range values {
		s.values[v] = struct{}{}
	}
}

// Remove removes 'val' from the set.
func (s Set[T]) Del(values ...T) {
	for _, v := range values {
		delete(s.values, v)
	}
}

// Has returns true only if 'val' is in the set.
func (s Set[T]) Has(val T) bool {
	_, ok := s.values[val]

	return ok
}

// Has returns true only if 'val' is in the set.
func (s Set[T]) HasAny(values ...T) bool {
	if len(values) == 0 {
		return true
	}

	for _, v := range values {
		if _, ok := s.values[v]; ok {
			return true
		}
	}

	return false
}

// Has returns true only if 'val' is in the set.
func (s Set[T]) HasAll(values ...T) bool {
	for _, v := range values {
		if _, ok := s.values[v]; !ok {
			return false
		}
	}

	return true
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s.values)
}

// Each calls 'fn' on every item in the set in no particular order.
func (s Set[T]) Each(fn func(key T)) {
	for k := range s.values {
		fn(k)
	}
}

// Values returns a slice of hashset elements.
func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s.values))

	for v := range s.values {
		values = append(values, v)
	}

	return values
}

// Sets operations.

// Union returns new set that contains all elements from both sets.
func (s Set[T]) Union(other Set[T]) Set[T] {
	newSet := make(map[T]struct{})

	for v := range s.values {
		newSet[v] = struct{}{}
	}

	for v := range other.values {
		newSet[v] = struct{}{}
	}

	return Set[T]{
		values: newSet,
	}
}

// Union returns new set that contains elements that are included in both sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	newSet := make(map[T]struct{})

	for v := range s.values {
		if _, ok := other.values[v]; ok {
			newSet[v] = struct{}{}
		}
	}

	return Set[T]{
		values: newSet,
	}
}

// Diff returns new set of values that are present in the first set, but not in the second.
func (s Set[T]) Diff(other Set[T]) Set[T] {
	newSet := make(map[T]struct{})

	for v := range s.values {
		if _, ok := other.values[v]; !ok {
			newSet[v] = struct{}{}
		}
	}

	return Set[T]{
		values: newSet,
	}
}

// IsSubset returns boolean whether fist set is subset of second.
func (s Set[T]) IsSubset(other Set[T]) bool {
	for v := range s.values {
		if _, ok := other.values[v]; !ok {
			return false
		}
	}

	return true
}
