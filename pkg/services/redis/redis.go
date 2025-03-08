package redisservices

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisService struct {
	redisURL string
	Client   *redis.Client
}

type Redis interface {
	BlackListToken(string, time.Duration) error
	IsTokenBlacklisted(string) bool
}

func NewRedisClient(redisURL string) Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL, // Redis default port
		Password: "",       // No password by default
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Failed to connect to Redis:", err)
		return nil
	}

	log.Println("Connected to Redis")
	return &RedisService{
		redisURL: redisURL,
		Client:   client,
	}
}

func (r *RedisService) BlackListToken(token string, exp time.Duration) error {
	return r.Client.Set(token, "blacklisted", exp).Err()
}

func (r *RedisService) IsTokenBlacklisted(token string) bool {
	_, err := r.Client.Get(token).Result()
	if err == redis.Nil {
		return false // Token is not blacklisted
	}
	if err != nil {
		fmt.Println("Redis error:", err)
		return false
	}
	return true // token is blacklisted
}
