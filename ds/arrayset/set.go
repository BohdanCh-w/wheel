// Package arrayset provides an implementation of a set using the built-in map.
package arrayset

import "fmt"

// New returns an empty arrayset.
func New[T comparable](values ...T) Set[T] {
	set := Set[T]{
		values: make([]T, 0, len(values)),
	}

	set.Add(values...)

	return set
}

// Set implements a arrayset, using an array (slice) as the underlying storage.
// This implementation is more efficient than the hashset for small sets.
type Set[T comparable] struct {
	values []T
}

func (s Set[T]) String() string {
	return fmt.Sprint(s.values)
}

// Add adds 'values' to the set.
func (s *Set[T]) Add(values ...T) {
	for _, v := range values {
		if !s.Has(v) {
			s.values = append(s.values, v)
		}
	}
}

// Remove removes 'val' from the set.
func (s *Set[T]) Del(values ...T) {
	if len(values) == 1 { // common case - use more efficient version
		for i, v := range s.values {
			if v == values[0] {
				s.values = append(s.values[:i], s.values[i+1:]...)

				return
			}
		}
	}

	var idx int

del:
	for i, v := range s.values {
		for _, d := range values {
			if v == d {
				values = values[1:]

				continue del
			}
		}

		s.values[idx] = s.values[i]
		idx++
	}

	s.values = s.values[:idx]
}

// Empty returns whether set has no elements.
func (s Set[T]) Empty() bool {
	return len(s.values) == 0
}

// Has returns true only if 'val' is in the set.
func (s Set[T]) Has(val T) bool {
	for _, v := range s.values {
		if v == val {
			return true
		}
	}

	return false
}

// Has returns true only if 'val' is in the set.
func (s Set[T]) HasAny(values ...T) bool {
	if len(values) == 0 {
		return true
	}

	for _, v := range values {
		if s.Has(v) {
			return true
		}
	}

	return false
}

// Has returns true only if 'val' is in the set.
func (s Set[T]) HasAll(values ...T) bool {
	for _, v := range values {
		if !s.Has(v) {
			return false
		}
	}

	return true
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s.values)
}

// Each calls 'fn' on every item in the set in order of adding.
func (s Set[T]) Each(fn func(key T)) {
	for _, v := range s.values {
		fn(v)
	}
}

// Map returns a new Set with applied function to each element.
func (s Set[T]) Map(fn func(key T) T) Set[T] {
	newSet := Set[T]{
		values: make([]T, 0, len(s.values)),
	}

	for _, v := range s.values {
		newSet.Add(fn(v))
	}

	return newSet
}

// Values returns a slice of arrayset elements.
func (s Set[T]) Values() []T {
	values := make([]T, len(s.values))
	copy(values, s.values)

	return values
}

// Clear clears all values.
func (s *Set[T]) Clear() {
	s.values = make([]T, 0)
}

// Returns copt of current set.
func (s Set[T]) Copy() Set[T] {
	return Set[T]{
		values: s.Values(),
	}
}

// Sets operations.

// Equal returns true if both sets have equal values.
func (s Set[T]) Equal(other Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}

	for _, v := range s.values {
		if !other.Has(v) {
			return false
		}
	}

	return true
}

// Union returns new set that contains all elements from both sets.
func (s Set[T]) Union(other Set[T]) Set[T] {
	newSet := Set[T]{
		values: make([]T, 0, s.Len()+other.Len()),
	}

	newSet.values = append(newSet.values, s.values...)

	for _, v := range other.values {
		newSet.Add(v)
	}

	return newSet
}

// Union returns new set that contains elements that are included in both sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	newSet := s.Copy()
	idx := 0

	for _, v := range newSet.values {
		if !other.Has(v) {
			continue
		}

		newSet.values[idx] = v
		idx++
	}

	newSet.values = newSet.values[:idx]

	return newSet
}

// Diff returns new set of values that are present in the first set, but not in the second.
func (s Set[T]) Diff(other Set[T]) Set[T] {
	newSet := s.Copy()
	idx := 0

	for _, v := range newSet.values {
		if other.Has(v) {
			continue
		}

		newSet.values[idx] = v
		idx++
	}

	newSet.values = newSet.values[:idx]

	return newSet
}

// SymDiff returns symmetric difference of two sets values.
func (s Set[T]) SymDiff(other Set[T]) Set[T] {
	return s.Union(other).Diff(s.Intersect(other))
}

// IsSubset returns boolean whether fist set is subset of second.
func (s Set[T]) IsSubset(other Set[T]) bool {
	for _, v := range s.values {
		if !other.Has(v) {
			return false
		}
	}

	return true
}
