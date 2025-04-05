package middleware

import (
	"cmp"
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
	ErrRateExceeded = wherr.Error("rate exceeded")

	clearRate       = time.Minute
	clearAfterTicks = 3
)

type rateLimiter interface {
	Allow() bool
}

type ClientIdentityFunc func(ctx context.Context, r *http.Request) string

type RateLimiterOpt func(*rateLimiterMiddleware)

func NewRateLimiter(
	logger logger.Logger,
	qps float64,
	opts ...RateLimiterOpt,
) (*rateLimiterMiddleware, context.CancelFunc) {
	shut := make(chan struct{})
	done := make(chan struct{})

	cancelFunc := func() {
		close(shut)
		<-done
	}

	limiter := &rateLimiterMiddleware{
		logger:         logger,
		qps:            qps,
		clients:        make(map[string]*rateClient),
		shut:           shut,
		done:           done,
		clientIDFunc:   defaultClientIdentityFunc,
		limiterFactory: defaultRateLimiterFactory(qps),
	}

	for _, opt := range opts {
		opt(limiter)
	}

	return limiter, cancelFunc
}

type rateLimiterMiddleware struct {
	logger       logger.Logger
	clientIDFunc ClientIdentityFunc
	qps          float64

	start          sync.Once
	shut           <-chan struct{}
	done           chan<- struct{}
	mu             sync.RWMutex
	clients        map[string]*rateClient
	limiterFactory func() rateLimiter
}

func (mid *rateLimiterMiddleware) Wrap(h api.Handler) api.Handler {
	mid.start.Do(func() {
		go mid.clean()
	})

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		log := logger.FromCtx(ctx, mid.logger)

		log.Debugf("start rate_limiter middleware")

		if id := mid.clientIDFunc(ctx, r); !mid.check(id) {
			log.Warnf("rate_limiter: %q client: too many requests", id)

			return web.NewError(http.StatusTooManyRequests, ErrRateExceeded)
		}

		return h(ctx, w, r)
	}
}

func (mid *rateLimiterMiddleware) clean() {
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

func (mid *rateLimiterMiddleware) check(id string) bool {
	mid.mu.RLock()
	client, ok := mid.clients[id]
	mid.mu.RUnlock()

	if !ok {
		client := &rateClient{
			limiter:        mid.limiterFactory(),
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
	limiter        rateLimiter
	lastConnection time.Time
}

func RateLimiterWithFactory(f func() rateLimiter) func(*rateLimiterMiddleware) {
	return func(rl *rateLimiterMiddleware) {
		rl.limiterFactory = f
	}
}

func RateLimiterWithIdentityFunc(f ClientIdentityFunc) func(*rateLimiterMiddleware) {
	return func(rl *rateLimiterMiddleware) {
		rl.clientIDFunc = f
	}
}

// default limiter burst size is half the qps
func defaultRateLimiterFactory(qps float64) func() rateLimiter {
	burstSize := cmp.Or(int(qps/2), 1)

	return func() rateLimiter {
		return rate.NewLimiter(rate.Limit(qps), burstSize)
	}
}

func defaultClientIdentityFunc(ctx context.Context, r *http.Request) string {
	return r.RemoteAddr
}
