package middleware

import (
	"errors"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/ratelimit"
	"net"
	"net/http"
)

var ErrRequestLimit = errors.New("you have reached the maximum number of requests or actions allowed within a certain time frame")

type RateLimiterMiddleware struct {
	service *ratelimit.RateLimiterService
}

func NewRateLimiterMiddleware(service *ratelimit.RateLimiterService) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{service: service}
}

func (rlm *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("API_KEY")
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		params := rlm.service.GetRateLimiterParams(apiKey, ip)

		if params.Visitor.Blocked {
			http.Error(w, ErrRequestLimit.Error(), http.StatusTooManyRequests)
			return
		}

		if !params.Limiter.Allow() {
			rlm.service.LockVisitor(ip)
			http.Error(w, ErrRequestLimit.Error(), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
