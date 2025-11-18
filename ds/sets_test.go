package ds_test

import (
	"math/rand/v2"
	"testing"

	"github.com/bohdanch-w/wheel/ds"
	"github.com/bohdanch-w/wheel/ds/arrayset"
	"github.com/bohdanch-w/wheel/ds/hashset"
	orderedset "github.com/bohdanch-w/wheel/ds/ordered-set"
)

func forEachSetType(b *testing.B, from []int, fn func(s ds.Set[int])) {
	b.Run("hashset", func(b *testing.B) {
		s := hashset.New[int](from...)

		for range b.N {
			fn(&s)
		}
	})

	b.Run("ordered", func(b *testing.B) {
		s := orderedset.New[int](from...)

		for range b.N {
			fn(&s)
		}
	})

	b.Run("arrayset", func(b *testing.B) {
		s := arrayset.New[int](from...)

		for range b.N {
			fn(&s)
		}
	})
}

func newInts(n int) []int {
	ints := make([]int, 0, n)

	for range n {
		ints = append(ints, rand.Int())
	}

	return ints
}

func BenchmarkSetsAddAndDel(b *testing.B) {
	execFunc := func(s ds.Set[int]) {
		n := rand.Int()

		s.Add(n)
		s.Del(n)
	}

	b.Run("xs", func(b *testing.B) {
		forEachSetType(b,
			newInts(5),
			execFunc,
		)
	})
	b.Run("s", func(b *testing.B) {
		forEachSetType(b,
			newInts(25),
			execFunc,
		)
	})
	b.Run("m", func(b *testing.B) {
		forEachSetType(b,
			newInts(100),
			execFunc,
		)
	})
	b.Run("l", func(b *testing.B) {
		forEachSetType(b,
			newInts(1000),
			execFunc,
		)
	})
	b.Run("xl", func(b *testing.B) {
		forEachSetType(b,
			newInts(10000),
			execFunc,
		)
	})
}

func BenchmarkSetsHas(b *testing.B) {
	execFunc := func(s ds.Set[int]) {
		n := rand.Int()

		s.Has(n)
	}

	b.Run("xs", func(b *testing.B) {
		forEachSetType(b,
			newInts(5),
			execFunc,
		)
	})
	b.Run("s", func(b *testing.B) {
		forEachSetType(b,
			newInts(25),
			execFunc,
		)
	})
	b.Run("m", func(b *testing.B) {
		forEachSetType(b,
			newInts(100),
			execFunc,
		)
	})
	b.Run("l", func(b *testing.B) {
		forEachSetType(b,
			newInts(1000),
			execFunc,
		)
	})
	b.Run("xl", func(b *testing.B) {
		forEachSetType(b,
			newInts(10000),
			execFunc,
		)
	})
}
