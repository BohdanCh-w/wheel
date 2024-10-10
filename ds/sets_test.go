package ds_test

import (
	"testing"

	"github.com/bohdanch-w/wheel/ds"
	"github.com/bohdanch-w/wheel/ds/arrayset"
	"github.com/bohdanch-w/wheel/ds/hashset"
	orderedset "github.com/bohdanch-w/wheel/ds/ordered-set"
)

func forEachSetType(b *testing.B, fn func(s ds.Set[int])) {
	b.Run("hashset", func(b *testing.B) {
		for range b.N {
			s := hashset.New[int]()

			fn(&s)
		}
	})

	b.Run("ordered", func(b *testing.B) {
		for range b.N {
			s := orderedset.New[int]()

			fn(&s)
		}
	})

	b.Run("arrayset", func(b *testing.B) {
		for range b.N {
			s := arrayset.New[int]()

			fn(&s)
		}
	})
}

func BenchmarkSetsAddSingle(b *testing.B) {
	b.Run("xs", func(b *testing.B) {
		forEachSetType(b, func(s ds.Set[int]) {
			for i := range 5 {
				s.Add(i)
			}
		})
	})
	b.Run("s", func(b *testing.B) {
		forEachSetType(b, func(s ds.Set[int]) {
			for i := range 25 {
				s.Add(i)
			}
		})
	})
	b.Run("m", func(b *testing.B) {
		forEachSetType(b, func(s ds.Set[int]) {
			for i := range 100 {
				s.Add(i)
			}
		})
	})
	b.Run("l", func(b *testing.B) {
		forEachSetType(b, func(s ds.Set[int]) {
			for i := range 1000 {
				s.Add(i)
			}
		})
	})
	b.Run("xl", func(b *testing.B) {
		forEachSetType(b, func(s ds.Set[int]) {
			for i := range 10000 {
				s.Add(i)
			}
		})
	})
}
