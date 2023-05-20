// Package symhashset provides an implementation of a symetric map using the built-in map.
// Validation on repeting keys is on user
package symmap

type Pair[L, R comparable] struct {
	Left  L
	Right R
}

func NewPairs[L, R comparable](pairs ...Pair[L, R]) SymMap[L, R] {
	set := SymMap[L, R]{
		left:  make(map[L]R),
		right: make(map[R]L),
	}

	for _, pair := range pairs {
		set.left[pair.Left] = pair.Right
		set.right[pair.Right] = pair.Left
	}

	return set
}

type SymMap[L, R comparable] struct {
	left  map[L]R
	right map[R]L
}

func (s SymMap[L, R]) Add(left L, right R) {
	s.left[left] = right
	s.right[right] = left
}

func (s SymMap[L, R]) Len() int {
	return len(s.left)
}

func (s SymMap[L, R]) Empty() bool {
	return len(s.left) == 0
}

func (s SymMap[L, R]) GetLeft(right R) L {
	return s.right[right]
}

func (s SymMap[L, R]) GetRight(left L) R {
	return s.left[left]
}
