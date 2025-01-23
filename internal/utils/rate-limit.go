package utils

import (
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"net/http"
	"time"
)

func getRateLimitMiddleware() *httplimit.Middleware {
	// Rate limiting
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   20,
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

var rateLimitMiddleware = getRateLimitMiddleware()

func RateLimitMiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rateLimitMiddleware.Handle(http.HandlerFunc(next)).ServeHTTP(w, r)
	}
}
