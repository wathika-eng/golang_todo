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
	r, err := client.Ping().Result()
	if err != nil {
		// exit if no redis client
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Printf("Connected to Redis: %v", r)
	return &RedisService{
		redisURL: redisURL,
		Client:   client,
	}
}

func (r *RedisService) BlackListToken(token string, exp time.Duration) error {
	return r.Client.Set(token, "blacklisted", exp).Err()
}

// panics if no redis client
func (r *RedisService) IsTokenBlacklisted(token string) bool {
	if r.Client == nil {
		fmt.Println("Redis client is not initialized")
		return false
	}

	_, err := r.Client.Get(token).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		fmt.Println("Redis error:", err)
		return false
	}
	return true
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
		return fmt.Errorf("error while marshalling the todo: %v", id)
	}
	status, err := r.Client.Set(id.String(), notesJson, 1*time.Hour).Result()
	if err != nil {
		return fmt.Errorf("unable to set cache key and value: %v", err)
	}
	log.Printf("cached: %v\n", status)
	return nil
}

func (r *RedisService) DeleteCache(id uuid.UUID) error {
	return r.Client.Del(id.String()).Err()
}
