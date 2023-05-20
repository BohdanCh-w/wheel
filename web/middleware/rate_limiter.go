package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	wherr "github.com/bohdanch-w/wheel/errors"
	"github.com/bohdanch-w/wheel/logger"
	"github.com/bohdanch-w/wheel/web"
	"github.com/bohdanch-w/wheel/web/api"
)

const (
	errRateExceeded = wherr.Error("rate exceeded")

	clearRate       = 1000
	clearAfterTicks = 3
	burstSize       = 1
)

type ClientIdentityFunc func(ctx context.Context, r *http.Request) string

type RateLimiter struct {
	Logger       logger.Logger
	ClientIDFunc ClientIdentityFunc
	QPS          float64
	Shut         <-chan struct{}

	mu      sync.Mutex
	clients map[string]*rateClient
}

func (mid *RateLimiter) Wrap(h api.Handler) api.Handler {
	mid.clients = make(map[string]*rateClient)

	go mid.clean()

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		log := logger.FromCtx(ctx, mid.Logger)

		log.Debugf("start rate_limiter middleware")

		if id := mid.ClientIDFunc(ctx, r); !mid.check(id) {
			log.Warnf("rate_limiter: %q client: too many requests", id)

			return web.NewError(http.StatusTooManyRequests, errRateExceeded)
		}

		return h(ctx, w, r)
	}
}

func (mid *RateLimiter) clean() {
	var (
		tickDuration = time.Duration(mid.QPS*float64(time.Second)) * clearRate
		ticker       = time.NewTicker(tickDuration)
	)

	defer ticker.Stop()

	for {
		select {
		case <-mid.Shut:
			mid.Logger.Infof("rate limiter done")

			return
		case <-ticker.C:
		}

		mid.mu.Lock()

		for key, client := range mid.clients {
			if time.Since(client.lastConnection) > clearAfterTicks {
				delete(mid.clients, key)
			}
		}

		mid.mu.Unlock()
	}
}

func (mid *RateLimiter) check(id string) bool {
	mid.mu.Lock()
	defer mid.mu.Unlock()

	if _, ok := mid.clients[id]; !ok {
		mid.clients[id] = &rateClient{
			limiter: rate.NewLimiter(rate.Limit(mid.QPS), burstSize),
		}
	}

	mid.clients[id].lastConnection = time.Now()

	return mid.clients[id].limiter.Allow()
}

type rateClient struct {
	limiter        *rate.Limiter
	lastConnection time.Time
}
