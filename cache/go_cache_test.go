package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("string/string", func(t *testing.T) {
		cache := NewGoCache[string, string](time.Second, 0, time.Minute)

		cache.Set("key", "value")
		value, found := cache.Get("key")
		require.True(t, found)
		require.Equal(t, "value", value)

		cache.Delete("key")
		value, found = cache.Get("key")
		require.False(t, found)
		require.Empty(t, value)
	})

	t.Run("int/struct", func(t *testing.T) {
		type Value struct {
			Something float64
		}

		v := Value{Something: 42.0}

		cache := NewGoCache[int, Value](time.Second, 0, time.Minute)

		cache.Set(42, v)
		value, found := cache.Get(42)
		require.True(t, found)
		require.Equal(t, v, value)

		cache.Delete(42)
		value, found = cache.Get(42)
		require.False(t, found)
		require.Empty(t, value)
	})
}
