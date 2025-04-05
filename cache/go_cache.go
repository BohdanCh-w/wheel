package cache

import (
	"math/rand/v2"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
)

var _ Cache[string, string] = (*GoCacheImpl[string, string])(nil)

func NewGoCache[K comparable, V any](
	expire time.Duration,
	expireJitter time.Duration,
	cleanUpInterval time.Duration,
) *GoCacheImpl[K, V] {
	var key K

	_, err := cast.ToStringE(key)
	if err != nil {
		panic("key must be stringifyable")
	}

	return &GoCacheImpl[K, V]{
		cache:        cache.New(expire, cleanUpInterval),
		expire:       expire,
		expireJitter: expireJitter,
	}
}

type GoCacheImpl[K comparable, V any] struct {
	cache        *cache.Cache
	expire       time.Duration
	expireJitter time.Duration
}

func (g *GoCacheImpl[K, V]) Get(key K) (zero V, found bool) {
	result, found := g.cache.Get(fastToString(key))
	if !found {
		return zero, false
	}

	return result.(V), true
}

func (g *GoCacheImpl[K, V]) Set(key K, value V) {
	g.cache.Set(fastToString(key), value, g.expiration())
}

func (g *GoCacheImpl[K, V]) Delete(key K) {
	g.cache.Delete(fastToString(key))
}

func (g *GoCacheImpl[K, V]) expiration() time.Duration {
	if g.expireJitter == 0 {
		return g.expire
	}

	return g.expire + time.Duration(rand.Float64()-0.5)*g.expireJitter
}

func fastToString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return cast.ToString(v)
	}
}
