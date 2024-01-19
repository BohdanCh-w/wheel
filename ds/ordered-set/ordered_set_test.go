package orderedset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	set "github.com/bohdanch-w/wheel/ordered-set"
)

func TestOrderedSetNewAndValues(t *testing.T) {
	emptySet := set.New[int]()
	fullSet := set.New(4, 6, 7)

	require.Equal(t, []int{}, emptySet.Values())
	require.Equal(t, []int{4, 6, 7}, fullSet.Values())
}

func TestOrderedSetAdd(t *testing.T) {
	emptySet := set.New[int]()
	require.Equal(t, []int{}, emptySet.Values())

	emptySet.Add(4)
	require.Equal(t, []int{4}, emptySet.Values())

	emptySet.Add(6, 7)
	require.Equal(t, []int{4, 6, 7}, emptySet.Values())

	emptySet.Add(4, 7)
	require.Equal(t, []int{4, 6, 7}, emptySet.Values())
}

func TestOrderedSetDel(t *testing.T) {
	fullSet := set.New(1, 2, 3, 4, 5, 6, 7)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, fullSet.Values())

	fullSet.Del(4)
	require.Equal(t, []int{1, 2, 3, 5, 6, 7}, fullSet.Values())

	fullSet.Del(6, 7)
	require.Equal(t, []int{1, 2, 3, 5}, fullSet.Values())

	fullSet.Del(4, 7)
	require.Equal(t, []int{1, 2, 3, 5}, fullSet.Values())

	fullSet.Del(1, 2, 3, 5, 8)
	require.Equal(t, []int{}, fullSet.Values())

	fullSet.Del(1, 2, 3)
	require.Equal(t, []int{}, fullSet.Values())
}

func TestOrderedSetEmpty(t *testing.T) {
	require.True(t, set.New[int]().Empty())
	require.False(t, set.New(1, 2, 3).Empty())
}

func TestOrderedSetLen(t *testing.T) {
	require.Equal(t, 0, set.New[int]().Len())
	require.Equal(t, 1, set.New(1).Len())
	require.Equal(t, 3, set.New(1, 2, 3).Len())
}

func TestOrderedSetClear(t *testing.T) {
	emptySet := set.New[int]()
	fullSet := set.New(4, 6, 7)

	emptySet.Clear()
	fullSet.Clear()

	require.Equal(t, 0, emptySet.Len())
	require.Equal(t, 0, fullSet.Len())
}

func TestOrderedSetCopy(t *testing.T) {
	set := set.New(1, 2, 3)
	clone := set.Copy()

	require.Equal(t, set.Values(), clone.Values())

	set.Add(4)

	require.Equal(t, []int{1, 2, 3, 4}, set.Values())
	require.Equal(t, []int{1, 2, 3}, clone.Values())
}

func TestOrderedSetEqual(t *testing.T) {
	setOne := set.New(1, 2, 3)
	setTwo := set.New(3, 4, 5)
	setThree := set.New(1, 2, 3, 4, 5)
	clone := setOne.Copy()

	require.True(t, setOne.Equal(clone))
	require.False(t, setOne.Equal(setTwo))
	require.False(t, setOne.Equal(setThree))
}

func TestOrderedSetEach(t *testing.T) {
	res := make([]int, 0)

	fn := func(v int) {
		res = append(res, v*2)
	}

	set.New[int]().Each(fn)
	require.Equal(t, []int{}, res)

	set.New(1).Each(fn)
	require.Equal(t, []int{2}, res)
	res = []int{} // reset

	set.New(1, 2, 3).Each(fn)
	require.Equal(t, []int{2, 4, 6}, res)
}

func TestOrderedSetMap(t *testing.T) {
	fn := func(v int) int {
		return v * 2
	}

	require.Equal(t, []int{}, set.New[int]().Map(fn).Values())
	require.Equal(t, []int{2}, set.New(1).Map(fn).Values())
	require.Equal(t, []int{2, 4, 6}, set.New(1, 2, 3).Map(fn).Values())
}

func TestOrderedSetHas(t *testing.T) {
	emptySet := set.New[int]()
	fullSet := set.New(4, 6, 7)

	require.False(t, emptySet.Has(4))
	require.True(t, fullSet.Has(4))

	require.True(t, fullSet.HasAny())
	require.True(t, fullSet.HasAny(4))
	require.False(t, fullSet.HasAny(5))
	require.True(t, fullSet.HasAny(2, 6, 8))
	require.False(t, fullSet.HasAny(2, 5, 8))

	require.True(t, fullSet.HasAll())
	require.True(t, fullSet.HasAll(4))
	require.False(t, fullSet.HasAll(5))
	require.True(t, fullSet.HasAll(4, 6, 7))
	require.False(t, fullSet.HasAll(4, 6, 8))
	require.False(t, fullSet.HasAll(2, 5, 8))
}

func TestOrderedSetUnion(t *testing.T) {
	require.Equal(t,
		[]int{},
		set.New[int]().Union(set.New[int]()).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3},
		set.New(1, 2, 3).Union(set.New[int]()).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3},
		set.New[int]().Union(set.New(1, 2, 3)).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3, 4, 5},
		set.New(1, 2, 3).Union(set.New(4, 5)).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3, 4, 5},
		set.New(1, 2, 3).Union(set.New(3, 4, 5)).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3},
		set.New(1, 2, 3).Union(set.New(2)).Values(),
	)
}

func TestOrderedSetIntersect(t *testing.T) {
	require.Equal(t,
		[]int{},
		set.New[int]().Intersect(set.New[int]()).Values(),
	)

	require.Equal(t,
		[]int{},
		set.New(1, 2, 3).Intersect(set.New[int]()).Values(),
	)

	require.Equal(t,
		[]int{},
		set.New[int]().Intersect(set.New(1, 2, 3)).Values(),
	)

	require.Equal(t,
		[]int{},
		set.New(1, 2, 3).Intersect(set.New(4, 5)).Values(),
	)

	require.Equal(t,
		[]int{3},
		set.New(1, 2, 3).Intersect(set.New(3, 4, 5)).Values(),
	)

	require.Equal(t,
		[]int{2},
		set.New(1, 2, 3).Intersect(set.New(2)).Values(),
	)

	require.Equal(t,
		[]int{2, 3, 4},
		set.New(1, 2, 3, 4, 5, 6).Intersect(set.New(2, 3, 4, 7, 8)).Values(),
	)
}

func TestOrderedSetDiff(t *testing.T) {
	require.Equal(t,
		[]int{},
		set.New[int]().Diff(set.New[int]()).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3},
		set.New(1, 2, 3).Diff(set.New[int]()).Values(),
	)

	require.Equal(t,
		[]int{},
		set.New[int]().Diff(set.New(1, 2, 3)).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3},
		set.New(1, 2, 3).Diff(set.New(4, 5)).Values(),
	)

	require.Equal(t,
		[]int{1, 2},
		set.New(1, 2, 3).Diff(set.New(3, 4, 5)).Values(),
	)

	require.Equal(t,
		[]int{1, 3},
		set.New(1, 2, 3).Diff(set.New(2)).Values(),
	)

	require.Equal(t,
		[]int{},
		set.New(2, 3).Diff(set.New(1, 2, 3, 4)).Values(),
	)
}

func TestOrderedSetSymDiff(t *testing.T) {
	require.Equal(t,
		[]int{},
		set.New[int]().SymDiff(set.New[int]()).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3},
		set.New(1, 2, 3).SymDiff(set.New[int]()).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3},
		set.New[int]().SymDiff(set.New(1, 2, 3)).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 3, 4, 5},
		set.New(1, 2, 3).SymDiff(set.New(4, 5)).Values(),
	)

	require.Equal(t,
		[]int{1, 2, 4, 5},
		set.New(1, 2, 3).SymDiff(set.New(3, 4, 5)).Values(),
	)

	require.Equal(t,
		[]int{1, 3},
		set.New(1, 2, 3).SymDiff(set.New(2)).Values(),
	)

	require.Equal(t,
		[]int{1, 4},
		set.New(2, 3).SymDiff(set.New(1, 2, 3, 4)).Values(),
	)
}

func TestOrderedSetIsSubset(t *testing.T) {
	require.True(t, set.New[int]().IsSubset(set.New[int]()))
	require.False(t, set.New(1, 2, 3).IsSubset(set.New[int]()))
	require.True(t, set.New[int]().IsSubset(set.New(1, 2, 3)))
	require.False(t, set.New(1, 2, 3).IsSubset(set.New(4, 5)))
	require.True(t, set.New(2, 3).IsSubset(set.New(1, 2, 3, 4)))
	require.False(t, set.New(1, 2, 3).IsSubset(set.New(2)))
}
