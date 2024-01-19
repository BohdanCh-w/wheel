package hashset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bohdanch-w/wheel/ds/hashset"
)

func TestHashSetNewAndValues(t *testing.T) {
	emptySet := hashset.New[int]()
	fullSet := hashset.New(4, 6, 7)

	require.ElementsMatch(t, []int{}, emptySet.Values())
	require.ElementsMatch(t, []int{4, 6, 7}, fullSet.Values())
}

func TestHashSetAdd(t *testing.T) {
	emptySet := hashset.New[int]()
	require.ElementsMatch(t, []int{}, emptySet.Values())

	emptySet.Add(4)
	require.ElementsMatch(t, []int{4}, emptySet.Values())

	emptySet.Add(6, 7)
	require.ElementsMatch(t, []int{4, 6, 7}, emptySet.Values())

	emptySet.Add(4, 7)
	require.ElementsMatch(t, []int{4, 6, 7}, emptySet.Values())
}

func TestHashSetDel(t *testing.T) {
	fullSet := hashset.New(1, 2, 3, 4, 5, 6, 7)
	require.ElementsMatch(t, []int{1, 2, 3, 4, 5, 6, 7}, fullSet.Values())

	fullSet.Del(4)
	require.ElementsMatch(t, []int{1, 2, 3, 5, 6, 7}, fullSet.Values())

	fullSet.Del(6, 7)
	require.ElementsMatch(t, []int{1, 2, 3, 5}, fullSet.Values())

	fullSet.Del(4, 7)
	require.ElementsMatch(t, []int{1, 2, 3, 5}, fullSet.Values())

	fullSet.Del(1, 2, 3, 5, 8)
	require.ElementsMatch(t, []int{}, fullSet.Values())

	fullSet.Del(1, 2, 3)
	require.ElementsMatch(t, []int{}, fullSet.Values())
}

func TestHashSetEmpty(t *testing.T) {
	require.True(t, hashset.New[int]().Empty())
	require.False(t, hashset.New(1, 2, 3).Empty())
}

func TestHashSetLen(t *testing.T) {
	require.Equal(t, 0, hashset.New[int]().Len())
	require.Equal(t, 1, hashset.New(1).Len())
	require.Equal(t, 3, hashset.New(1, 2, 3).Len())
}

func TestHashSetClear(t *testing.T) {
	emptySet := hashset.New[int]()
	fullSet := hashset.New(4, 6, 7)

	emptySet.Clear()
	fullSet.Clear()

	require.Equal(t, 0, emptySet.Len())
	require.Equal(t, 0, fullSet.Len())
}

func TestHashSetCopy(t *testing.T) {
	set := hashset.New(1, 2, 3)
	clone := set.Copy()

	require.ElementsMatch(t, set.Values(), clone.Values())

	set.Add(4)

	require.ElementsMatch(t, []int{1, 2, 3, 4}, set.Values())
	require.ElementsMatch(t, []int{1, 2, 3}, clone.Values())
}

func TestHashSetEqual(t *testing.T) {
	setOne := hashset.New(1, 2, 3)
	setTwo := hashset.New(3, 4, 5)
	setThree := hashset.New(1, 2, 3, 4, 5)
	clone := setOne.Copy()

	require.True(t, setOne.Equal(clone))
	require.False(t, setOne.Equal(setTwo))
	require.False(t, setOne.Equal(setThree))
}

func TestHashSetEach(t *testing.T) {
	res := make([]int, 0)

	fn := func(v int) {
		res = append(res, v*2)
	}

	hashset.New[int]().Each(fn)
	require.ElementsMatch(t, []int{}, res)

	hashset.New(1).Each(fn)
	require.ElementsMatch(t, []int{2}, res)
	res = []int{} // reset

	hashset.New(1, 2, 3).Each(fn)
	require.ElementsMatch(t, []int{2, 4, 6}, res)
}

func TestHashSetMap(t *testing.T) {
	fn := func(v int) int {
		return v * 2
	}

	require.ElementsMatch(t, []int{}, hashset.New[int]().Map(fn).Values())
	require.ElementsMatch(t, []int{2}, hashset.New(1).Map(fn).Values())
	require.ElementsMatch(t, []int{2, 4, 6}, hashset.New(1, 2, 3).Map(fn).Values())
}

func TestHashSetHas(t *testing.T) {
	emptySet := hashset.New[int]()
	fullSet := hashset.New(4, 6, 7)

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

func TestHashSetUnion(t *testing.T) {
	require.ElementsMatch(t,
		[]int{},
		hashset.New[int]().Union(hashset.New[int]()).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3},
		hashset.New(1, 2, 3).Union(hashset.New[int]()).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3},
		hashset.New[int]().Union(hashset.New(1, 2, 3)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3, 4, 5},
		hashset.New(1, 2, 3).Union(hashset.New(4, 5)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3, 4, 5},
		hashset.New(1, 2, 3).Union(hashset.New(3, 4, 5)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3},
		hashset.New(1, 2, 3).Union(hashset.New(2)).Values(),
	)
}

func TestHashSetIntersect(t *testing.T) {
	require.ElementsMatch(t,
		[]int{},
		hashset.New[int]().Intersect(hashset.New[int]()).Values(),
	)

	require.ElementsMatch(t,
		[]int{},
		hashset.New(1, 2, 3).Intersect(hashset.New[int]()).Values(),
	)

	require.ElementsMatch(t,
		[]int{},
		hashset.New[int]().Intersect(hashset.New(1, 2, 3)).Values(),
	)

	require.ElementsMatch(t,
		[]int{},
		hashset.New(1, 2, 3).Intersect(hashset.New(4, 5)).Values(),
	)

	require.ElementsMatch(t,
		[]int{3},
		hashset.New(1, 2, 3).Intersect(hashset.New(3, 4, 5)).Values(),
	)

	require.ElementsMatch(t,
		[]int{2},
		hashset.New(1, 2, 3).Intersect(hashset.New(2)).Values(),
	)

	require.ElementsMatch(t,
		[]int{2, 3, 4},
		hashset.New(1, 2, 3, 4, 5, 6).Intersect(hashset.New(2, 3, 4, 7, 8)).Values(),
	)
}

func TestHashSetDiff(t *testing.T) {
	require.ElementsMatch(t,
		[]int{},
		hashset.New[int]().Diff(hashset.New[int]()).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3},
		hashset.New(1, 2, 3).Diff(hashset.New[int]()).Values(),
	)

	require.ElementsMatch(t,
		[]int{},
		hashset.New[int]().Diff(hashset.New(1, 2, 3)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3},
		hashset.New(1, 2, 3).Diff(hashset.New(4, 5)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2},
		hashset.New(1, 2, 3).Diff(hashset.New(3, 4, 5)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 3},
		hashset.New(1, 2, 3).Diff(hashset.New(2)).Values(),
	)

	require.ElementsMatch(t,
		[]int{},
		hashset.New(2, 3).Diff(hashset.New(1, 2, 3, 4)).Values(),
	)
}

func TestHashSetSymDiff(t *testing.T) {
	require.ElementsMatch(t,
		[]int{},
		hashset.New[int]().SymDiff(hashset.New[int]()).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3},
		hashset.New(1, 2, 3).SymDiff(hashset.New[int]()).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3},
		hashset.New[int]().SymDiff(hashset.New(1, 2, 3)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 3, 4, 5},
		hashset.New(1, 2, 3).SymDiff(hashset.New(4, 5)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 2, 4, 5},
		hashset.New(1, 2, 3).SymDiff(hashset.New(3, 4, 5)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 3},
		hashset.New(1, 2, 3).SymDiff(hashset.New(2)).Values(),
	)

	require.ElementsMatch(t,
		[]int{1, 4},
		hashset.New(2, 3).SymDiff(hashset.New(1, 2, 3, 4)).Values(),
	)
}

func TestHashSetIsSubset(t *testing.T) {
	require.True(t, hashset.New[int]().IsSubset(hashset.New[int]()))
	require.False(t, hashset.New(1, 2, 3).IsSubset(hashset.New[int]()))
	require.True(t, hashset.New[int]().IsSubset(hashset.New(1, 2, 3)))
	require.False(t, hashset.New(1, 2, 3).IsSubset(hashset.New(4, 5)))
	require.True(t, hashset.New(2, 3).IsSubset(hashset.New(1, 2, 3, 4)))
	require.False(t, hashset.New(1, 2, 3).IsSubset(hashset.New(2)))
}
