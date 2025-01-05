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

	clearRate       = time.Minute
	clearAfterTicks = 3
	burstSize       = 1
)

type ClientIdentityFunc func(ctx context.Context, r *http.Request) string

func ClientIdentityIPFunc(ctx context.Context, r *http.Request) string {
	return r.RemoteAddr
}

func NewRateLimiter(
	logger logger.Logger,
	clientIDFunc ClientIdentityFunc,
	qps float64,
) (*rateLimiter, context.CancelFunc) {
	shut := make(chan struct{})
	done := make(chan struct{})

	cancelFunc := func() {
		close(shut)
		<-done
	}

	limiter := &rateLimiter{
		logger:       logger,
		clientIDFunc: clientIDFunc,
		qps:          qps,
		clients:      make(map[string]*rateClient),
		shut:         shut,
		done:         done,
	}

	return limiter, cancelFunc
}

type rateLimiter struct {
	logger       logger.Logger
	clientIDFunc ClientIdentityFunc
	qps          float64

	start   sync.Once
	shut    <-chan struct{}
	done    chan<- struct{}
	mu      sync.RWMutex
	clients map[string]*rateClient
}

func (mid *rateLimiter) Wrap(h api.Handler) api.Handler {
	mid.start.Do(func() {
		go mid.clean()
	})

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		log := logger.FromCtx(ctx, mid.logger)

		log.Debugf("start rate_limiter middleware")

		if id := mid.clientIDFunc(ctx, r); !mid.check(id) {
			log.Warnf("rate_limiter: %q client: too many requests", id)

			return web.NewError(http.StatusTooManyRequests, errRateExceeded)
		}

		return h(ctx, w, r)
	}
}

func (mid *rateLimiter) clean() {
	var (
		tickDuration = max(clearRate/time.Duration(mid.qps), time.Second)
		ticker       = time.NewTicker(tickDuration)
	)

	defer ticker.Stop()
	defer close(mid.done)

	for {
		select {
		case <-mid.shut:
			mid.logger.Infof("rate limiter done")

			return
		case <-ticker.C:
		}

		mid.mu.Lock()

		for key, client := range mid.clients {
			if time.Since(client.lastConnection) > clearAfterTicks*tickDuration {
				delete(mid.clients, key)
			}
		}

		mid.mu.Unlock()
	}
}

func (mid *rateLimiter) check(id string) bool {
	mid.mu.RLock()
	client, ok := mid.clients[id]
	mid.mu.RUnlock()

	if !ok {
		client := &rateClient{
			limiter:        rate.NewLimiter(rate.Limit(mid.qps), burstSize),
			lastConnection: time.Now(),
		}

		mid.mu.Lock()
		mid.clients[id] = client
		mid.mu.Unlock()

		return client.limiter.Allow()
	}

	client.lastConnection = time.Now()

	return mid.clients[id].limiter.Allow()
}

type rateClient struct {
	limiter        *rate.Limiter
	lastConnection time.Time
}
