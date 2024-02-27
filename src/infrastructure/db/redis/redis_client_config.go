package redis

import (
	"fmt"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/env"
	"github.com/go-redis/redis"
	"log"
)

func Initialize(cfg env.EnvConfig) *redis.Client {
	options := &redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	client := redis.NewClient(options)

	_, err := client.Ping().Result()

	if err != nil {
		log.Fatalf("Connection to Redis failed: %v", err)
		return nil
	}

	log.Println(fmt.Sprintf("Redis Initialized: %+v", options))
	return client
}
