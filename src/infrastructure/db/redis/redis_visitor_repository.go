package redis

import (
	"encoding/json"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/domain"
	"github.com/go-redis/redis"
	"log"
)

type VisitorRepository struct {
	client *redis.Client
}

func NewRedisVisitorRepository(client *redis.Client) *VisitorRepository {
	return &VisitorRepository{client: client}
}

func (r VisitorRepository) SaveOrUpdate(entity domain.VisitorEntity) {
	jsonEntity, err := json.Marshal(entity)
	if err != nil {
		log.Printf("Error when trying to marshal entity: %v\n", err)
		return
	}
	err = r.client.Set(entity.IpAddress, jsonEntity, 0).Err()
	if err != nil {
		log.Printf("Error when trying to save data to Redis: %v\n", err)
	}
}

func (r VisitorRepository) Find(ipAddress string) *domain.VisitorEntity {
	result, err := r.client.Get(ipAddress).Result()
	if err != nil {
		log.Printf("Error when trying to get data from Redis: %v\n", err)
		return nil
	}

	var entity domain.VisitorEntity
	err = json.Unmarshal([]byte(result), &entity)
	if err != nil {
		log.Printf("Error when trying to unmarshal data: %v\n", err)
		return nil
	}

	return &entity
}
