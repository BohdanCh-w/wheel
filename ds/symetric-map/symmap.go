// Package symhashset provides an implementation of a symetric map using the built-in map.
// Validation on repeting keys is on user
package symmap

import (
	"fmt"

	"github.com/bohdanch-w/wheel/errors"
)

type Pair[L, R comparable] struct {
	L L
	R R
}

func New[L, R comparable](pairs ...any) *SymMap[L, R] {
	mp, err := NewValues[L, R](pairs...)
	if err != nil {
		panic(err)
	}

	return mp
}

func NewPairs[L, R comparable](pairs ...Pair[L, R]) (*SymMap[L, R], error) {
	mp := &SymMap[L, R]{
		left:  make(map[L]R),
		right: make(map[R]L),
	}

	for _, pair := range pairs {
		if err := mp.Add(pair.L, pair.R); err != nil {
			return nil, err
		}
	}

	return mp, nil
}

func NewValues[L, R comparable](pairs ...any) (*SymMap[L, R], error) {
	const (
		errInvalidParamNumber = errors.Error("invalid number of params")
		errInvalidPair        = errors.Error("invalid pair")
	)

	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("%w: %d", errInvalidParamNumber, len(pairs))
	}

	mp := &SymMap[L, R]{
		left:  make(map[L]R),
		right: make(map[R]L),
	}

	for i := 0; i < len(pairs); i += 2 {
		left, lOk := pairs[i].(L)
		right, rOk := pairs[i+1].(R)

		if !lOk || !rOk {
			return nil, fmt.Errorf(
				"%w: received %v[%T] - %v[%T], expected [%T] - [%T]",
				errInvalidPair,
				pairs[i], pairs[i],
				pairs[i+1], pairs[i+1],
				left, right,
			)
		}

		if err := mp.Add(left, right); err != nil {
			return nil, err
		}
	}

	return mp, nil
}

// SymMap is a standard golang map that can be used to map values in both directions.
type SymMap[L, R comparable] struct {
	left  map[L]R
	right map[R]L
}

func (s SymMap[L, R]) Add(left L, right R) error {
	const errAlreadyMapped = errors.Error("duplicate value")

	if right, ok := s.left[left]; ok {
		return fmt.Errorf("%w: %v - %v", errAlreadyMapped, left, right)
	}

	if left, ok := s.right[right]; ok {
		return fmt.Errorf("%w: %v - %v", errAlreadyMapped, left, right)
	}

	s.left[left] = right
	s.right[right] = left

	return nil
}

func (s SymMap[L, R]) Len() int {
	return len(s.left)
}

func (s SymMap[L, R]) Empty() bool {
	return len(s.left) == 0
}

func (s SymMap[L, R]) Left(right R) (L, bool) {
	left, ok := s.right[right]

	return left, ok
}

func (s SymMap[L, R]) Right(left L) (R, bool) {
	right, ok := s.left[left]

	return right, ok
}

func (s SymMap[L, R]) GetLeft(right R) L {
	return s.right[right]
}

func (s SymMap[L, R]) GetRight(left L) R {
	return s.left[left]
}

func (s SymMap[L, R]) Values() []Pair[L, R] {
	pairs := make([]Pair[L, R], 0, s.Len())

	for left, right := range s.left {
		pairs = append(pairs, Pair[L, R]{L: left, R: right})
	}

	return pairs
}

func (s SymMap[L, R]) Map() map[L]R {
	pairs := make(map[L]R, s.Len())

	for left, right := range s.left {
		pairs[left] = right
	}

	return pairs
}
