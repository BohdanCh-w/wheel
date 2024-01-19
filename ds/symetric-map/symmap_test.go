package symmap_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	symmap "github.com/bohdanch-w/wheel/ds/symetric-map"
)

func TestSymetricMapNew_Success(t *testing.T) {
	emptyMap, err := symmap.NewValues[string, int]()
	require.NoError(t, err)

	mp, err := symmap.NewValues[string, int]("test1", 1, "test2", 2)
	require.NoError(t, err)

	require.Equal(t, map[string]int{}, emptyMap.Map())
	require.Equal(t, map[string]int{"test1": 1, "test2": 2}, mp.Map())
}

func TestSymetricMapNew_FailsNumberOfParams(t *testing.T) {
	mp, err := symmap.NewValues[string, int]("test1", 1, "test2", 2, "test3")
	require.Empty(t, mp)
	require.EqualError(t, err, "invalid number of params: 5")
}

func TestSymetricMapNew_FailsType(t *testing.T) {
	mp, err := symmap.NewValues[string, int]("test1", 1, 2.0, 2)
	require.Empty(t, mp)
	require.EqualError(t, err, "invalid pair: received 2[float64] - 2[int], expected [string] - [int]")
}

func TestSymetricMapNew_FailsDuplicate(t *testing.T) {
	mp, err := symmap.NewValues[string, int]("test1", 1, "test1", 2)
	require.Empty(t, mp)
	require.EqualError(t, err, "duplicate value: test1 - 1")

	mp, err = symmap.NewValues[string, int]("test1", 1, "test2", 1)
	require.Empty(t, mp)
	require.EqualError(t, err, "duplicate value: test1 - 1")

	mp, err = symmap.NewValues[string, int]("test1", 1, "test1", 1)
	require.Empty(t, mp)
	require.EqualError(t, err, "duplicate value: test1 - 1")
}

func TestSymetricMapNewPairs_Success(t *testing.T) {
	mp, err := symmap.NewPairs(symmap.Pair[string, int]{"test1", 1}, symmap.Pair[string, int]{"test2", 2})
	require.NoError(t, err)

	require.Equal(t, map[string]int{"test1": 1, "test2": 2}, mp.Map())
}

func TestSymetricMapNewPairs_FailsDuplicate(t *testing.T) {
	mp, err := symmap.NewPairs(symmap.Pair[string, int]{"test1", 1}, symmap.Pair[string, int]{"test1", 2})
	require.Empty(t, mp)
	require.EqualError(t, err, "duplicate value: test1 - 1")

	mp, err = symmap.NewPairs(symmap.Pair[string, int]{"test1", 1}, symmap.Pair[string, int]{"test2", 1})
	require.Empty(t, mp)
	require.EqualError(t, err, "duplicate value: test1 - 1")

	mp, err = symmap.NewPairs(symmap.Pair[string, int]{"test1", 1}, symmap.Pair[string, int]{"test1", 1})
	require.Empty(t, mp)
	require.EqualError(t, err, "duplicate value: test1 - 1")
}

func TestSymetricMapMust_Success(t *testing.T) {
	require.NotPanics(t, func() {
		mp := symmap.New[string, int]("test1", 1, "test2", 2)

		require.Equal(t, map[string]int{"test1": 1, "test2": 2}, mp.Map())
	})
}

func TestSymetricMapMust_Fail(t *testing.T) {
	require.PanicsWithError(t, "invalid number of params: 5", func() {
		symmap.New[string, int]("test1", 1, "test2", 2, "test3")
	})
}

func TestSymetricMapLen(t *testing.T) {
	emptyMap := symmap.New[string, int]()
	mp := symmap.New[string, int]("test1", 1)

	require.Equal(t, 0, emptyMap.Len())
	require.Equal(t, 1, mp.Len())
}

func TestSymetricMapEmpty(t *testing.T) {
	emptyMap := symmap.New[string, int]()
	mp := symmap.New[string, int]("test1", 1)

	require.True(t, emptyMap.Empty())
	require.False(t, mp.Empty())
}

func TestSymetricMapAdd_Success(t *testing.T) {
	mp := symmap.New[string, int]()

	require.Equal(t, 0, mp.Len())

	err := mp.Add("test1", 1)
	require.NoError(t, err)

	err = mp.Add("test2", 2)
	require.NoError(t, err)

	require.Equal(t, map[string]int{"test1": 1, "test2": 2}, mp.Map())
}

func TestSymetricMapAdd_Fail(t *testing.T) {
	mp := symmap.New[string, int]()

	require.Equal(t, 0, mp.Len())

	err := mp.Add("test1", 1)
	require.NoError(t, err)

	err = mp.Add("test1", 1)
	require.EqualError(t, err, "duplicate value: test1 - 1")
}

func TestSymetricMapValues(t *testing.T) {
	emptyMap := symmap.New[string, int]()
	mp := symmap.New[string, int]("test1", 1, "test2", 2)

	require.ElementsMatch(t, []symmap.Pair[string, int]{}, emptyMap.Values())
	require.ElementsMatch(t, []symmap.Pair[string, int]{{L: "test1", R: 1}, {L: "test2", R: 2}}, mp.Values())
}

func TestSymetricMapRight(t *testing.T) {
	mp := symmap.New[string, int]("test1", 1)

	right, ok := mp.Right("test1")
	require.True(t, ok)
	require.Equal(t, 1, right)

	right, ok = mp.Right("test2")
	require.False(t, ok)
	require.Empty(t, right)
}

func TestSymetricMapLeft(t *testing.T) {
	mp := symmap.New[string, int]("test1", 1)

	left, ok := mp.Left(1)
	require.True(t, ok)
	require.Equal(t, "test1", left)

	left, ok = mp.Left(2)
	require.False(t, ok)
	require.Empty(t, left)
}

func TestSymetricMapGetRight(t *testing.T) {
	mp := symmap.New[string, int]("test1", 1)

	require.Equal(t, 1, mp.GetRight("test1"))
	require.Equal(t, 0, mp.GetRight("test2"))
}

func TestSymetricMapGetLeft(t *testing.T) {
	mp := symmap.New[string, int]("test1", 1)

	require.Equal(t, "test1", mp.GetLeft(1))
	require.Equal(t, "", mp.GetLeft(2))
}
