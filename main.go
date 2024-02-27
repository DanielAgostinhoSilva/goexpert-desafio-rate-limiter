package main

import (
	"fmt"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/db/redis"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/env"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/ratelimit"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/webserver/middleware"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func main() {
	cfg := env.LoadConfig("./.env")
	redisClient := redis.Initialize(*cfg)
	repository := redis.NewRedisVisitorRepository(redisClient)
	service := ratelimit.NewRateLimiterService(*cfg, repository)
	defer service.CleanupVisitors()
	middleware := middleware.NewRateLimiterMiddleware(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":"+cfg.WebServerPort, middleware.Handler(mux))
}
