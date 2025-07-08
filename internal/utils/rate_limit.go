package utils

import (
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"net/http"
	"sync"
	"time"
)

func getRateLimitMiddleware(cfg *AppConfig) *httplimit.Middleware {
	// Rate limiting
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   cfg.RateLimitTokens,
		Interval: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	middleware, err := httplimit.NewMiddleware(store, httplimit.IPKeyFunc())
	if err != nil {
		panic(err)
	}

	return middleware
}

var rateLimitMiddleware *httplimit.Middleware
var rateLimitMiddlewareOnce sync.Once

func RateLimitMiddlewareFunc(cfg *AppConfig, next http.HandlerFunc) http.HandlerFunc {
	rateLimitMiddlewareOnce.Do(func() {
		rateLimitMiddleware = getRateLimitMiddleware(cfg)
	})

	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.RateLimitTokens > 0 {
			rateLimitMiddleware.Handle(next).ServeHTTP(w, r)
		} else {
			next(w, r)
		}
	}
}
