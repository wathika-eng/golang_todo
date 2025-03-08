package redisservices

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisClient() *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis default port
		Password: "",               // No password by default
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Failed to connect to Redis:", err)
		return nil
	}

	log.Println("Connected to Redis")
	return &RedisService{Client: client}
}

func (r *RedisService) BlackListToken(token string, exp time.Duration) error {
	return r.Client.Set("blacklisted", token, exp).Err()
}

func (r *RedisService) IsTokenBlacklisted(token string) bool {
	_, err := r.Client.Get(token).Result()
	return err == nil
}
