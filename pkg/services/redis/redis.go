package redisservices

import (
	"encoding/json"
	"fmt"
	"golang_todo/pkg/types"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type RedisService struct {
	redisURL string
	Client   *redis.Client
}

type Redis interface {
	BlackListToken(string, time.Duration) error
	IsTokenBlacklisted(string) bool
	FetchFromCache(id uuid.UUID) ([]types.Note, error)
	CacheTodo(todo interface{}, id uuid.UUID) error
	DeleteCache(id uuid.UUID) error
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

func (r *RedisService) FetchFromCache(id uuid.UUID) ([]types.Note, error) {
	cachedTodo, err := r.Client.Get(id.String()).Result()
	if err == nil {
		var todo []types.Note
		if json.Unmarshal([]byte(cachedTodo), &todo) == nil {
			return todo, nil
		}
	}
	return nil, err
}

func (r *RedisService) CacheTodo(todo interface{}, id uuid.UUID) error {
	notesJson, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("could not cache the todo: %v", id)
	}
	r.Client.Set(id.String(), notesJson, 1*time.Hour)
	return redis.Nil
}

func (r *RedisService) DeleteCache(id uuid.UUID) error {
	return r.Client.Del(id.String()).Err()
}
