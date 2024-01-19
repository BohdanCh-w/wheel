// Package orderedset provides an implementation of a set using the built-in map.
package orderedset

import "slices"

// New returns an empty orderedset.
func New[T comparable](values ...T) Set[T] {
	set := Set[T]{
		values: make(map[T]int),
	}

	set.Add(values...)

	return set
}

// Set implements a orderedset, using the hashmap as the underlying storage.
type Set[T comparable] struct {
	values map[T]int
	count  int
}

// Add adds 'values' to the set.
func (s *Set[T]) Add(values ...T) {
	for _, v := range values {
		if _, ok := s.values[v]; !ok {
			s.values[v] = s.count
			s.count++
		}
	}
}

// Remove removes 'val' from the set.
func (s Set[T]) Del(values ...T) {
	for _, v := range values {
		delete(s.values, v)
	}
}

// Empty retruns whether set has no elements.
func (s Set[T]) Empty() bool {
	return len(s.values) == 0
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

// Map returns a new Set with applied function to each element.
func (s Set[T]) Map(fn func(key T) T) Set[T] {
	newSet := make(map[T]int)

	for k, ord := range s.values {
		newSet[fn(k)] = ord
	}

	return Set[T]{
		values: newSet,
		count:  s.count,
	}
}

// Values returns a slice of orderedset elements.
func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s.values))

	for v := range s.values {
		values = append(values, v)
	}

	slices.SortFunc(values, func(v1, v2 T) int {
		return s.values[v1] - s.values[v2]
	})

	return values
}

// Clear clears all values.
func (s *Set[T]) Clear() {
	s.values = make(map[T]int)
}

// Returns copt of current set.
func (s Set[T]) Copy() Set[T] {
	return New(s.Values()...)
}

// Sets operations.

// Equal returns true if both sets have equal values.
func (s Set[T]) Equal(other Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}

	for v := range s.values {
		if _, ok := other.values[v]; !ok {
			return false
		}
	}

	return true
}

// Union returns new set that contains all elements from both sets.
func (s Set[T]) Union(other Set[T]) Set[T] {
	newSet := make(map[T]int)

	maxOrd := 0

	for v, ord := range s.values {
		newSet[v] = ord

		if ord > maxOrd {
			maxOrd = ord
		}
	}

	maxOrd++

	secMaxOrd := maxOrd

	for v, ord := range other.values {
		if _, ok := newSet[v]; ok {
			continue
		}

		newSet[v] = maxOrd + ord

		if secOrd := maxOrd + ord; secOrd > secMaxOrd {
			secMaxOrd = secOrd
		}
	}

	secMaxOrd++

	return Set[T]{
		values: newSet,
		count:  secMaxOrd,
	}
}

// Union returns new set that contains elements that are included in both sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	newSet := make(map[T]int)

	for v, ord := range s.values {
		if _, ok := other.values[v]; ok {
			newSet[v] = ord
		}
	}

	return Set[T]{
		values: newSet,
	}
}

// Diff returns new set of values that are present in the first set, but not in the second.
func (s Set[T]) Diff(other Set[T]) Set[T] {
	newSet := make(map[T]int)

	for v, ord := range s.values {
		if _, ok := other.values[v]; !ok {
			newSet[v] = ord
		}
	}

	return Set[T]{
		values: newSet,
		count:  s.count,
	}
}

// SymDiff returns symetric difference of two sets values.
func (s Set[T]) SymDiff(other Set[T]) Set[T] {
	newSet := s.Union(other)

	for v := range s.values {
		if _, ok := other.values[v]; ok {
			delete(newSet.values, v)
		}
	}

	return newSet
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
