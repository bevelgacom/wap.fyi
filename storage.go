package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// ChallengeStorage interface defines the methods for storing and retrieving challenges
type ChallengeStorage interface {
	Store(challenge string, solved bool) error
	Get(challenge string) (bool, bool, error) // returns (solved, exists, error)
	StoreURL(path string, fullURL string) error
	GetURL(path string) (string, bool, error) // returns (fullURL, exists, error)
	Close() error
}

// RedisStorage implements ChallengeStorage using Redis/Valkey
type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStorage creates a new Redis storage instance
func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// Store stores a challenge with its solved status in Redis
func (r *RedisStorage) Store(challenge string, solved bool) error {
	key := fmt.Sprintf("challenge:%s", challenge)
	value := "0"
	if solved {
		value = "1"
	}

	// Set with 24-hour expiration
	err := r.client.Set(r.ctx, key, value, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store challenge in Redis: %w", err)
	}

	return nil
}

// Get retrieves a challenge's solved status from Redis
func (r *RedisStorage) Get(challenge string) (bool, bool, error) {
	key := fmt.Sprintf("challenge:%s", challenge)

	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return false, false, nil // Key doesn't exist
	} else if err != nil {
		return false, false, fmt.Errorf("failed to get challenge from Redis: %w", err)
	}

	solved := val == "1"
	return solved, true, nil
}

// Close closes the Redis connection
func (r *RedisStorage) Close() error {
	return r.client.Close()
}

// StoreURL stores a URL mapping in Redis
func (r *RedisStorage) StoreURL(path string, fullURL string) error {
	key := fmt.Sprintf("url:%s", path)

	// Set with 24-hour expiration (same as challenges)
	err := r.client.Set(r.ctx, key, fullURL, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store URL in Redis: %w", err)
	}

	return nil
}

// GetURL retrieves a URL mapping from Redis
func (r *RedisStorage) GetURL(path string) (string, bool, error) {
	key := fmt.Sprintf("url:%s", path)

	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil // Key doesn't exist
	} else if err != nil {
		return "", false, fmt.Errorf("failed to get URL from Redis: %w", err)
	}

	return val, true, nil
}

// LocalMapStorage implements ChallengeStorage using an in-memory map
type LocalMapStorage struct {
	challenges map[string]bool
	urls       map[string]string
	mu         sync.RWMutex
}

// NewLocalMapStorage creates a new local map storage instance
func NewLocalMapStorage() *LocalMapStorage {
	return &LocalMapStorage{
		challenges: make(map[string]bool),
		urls:       make(map[string]string),
	}
}

// Store stores a challenge with its solved status in the local map
func (l *LocalMapStorage) Store(challenge string, solved bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.challenges[challenge] = solved
	return nil
}

// Get retrieves a challenge's solved status from the local map
func (l *LocalMapStorage) Get(challenge string) (bool, bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	solved, exists := l.challenges[challenge]
	return solved, exists, nil
}

// Close is a no-op for local map storage
func (l *LocalMapStorage) Close() error {
	return nil
}

// StoreURL stores a URL mapping in the local map
func (l *LocalMapStorage) StoreURL(path string, fullURL string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.urls[path] = fullURL
	return nil
}

// GetURL retrieves a URL mapping from the local map
func (l *LocalMapStorage) GetURL(path string) (string, bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	fullURL, exists := l.urls[path]
	return fullURL, exists, nil
}

// NewChallengeStorage creates a new challenge storage instance based on environment
func NewChallengeStorage() ChallengeStorage {
	// Check if we should use Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	useRedis := os.Getenv("USE_REDIS")

	// Use Redis if explicitly enabled or if running in production
	if useRedis == "true" || os.Getenv("ENV") == "production" {
		redisStorage, err := NewRedisStorage(redisAddr, redisPassword, 0)
		if err != nil {
			log.Printf("Failed to initialize Redis storage: %v. Falling back to local storage.", err)
			return NewLocalMapStorage()
		}
		log.Println("Using Redis storage for challenges")
		return redisStorage
	}

	log.Println("Using local map storage for challenges")
	return NewLocalMapStorage()
}
