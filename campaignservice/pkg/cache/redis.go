package cache

import (
	"context"
	"fmt"
	"log"
	"specommerce/campaignservice/pkg/service_config"
	"specommerce/campaignservice/pkg/shutdown"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(config service_config.RedisConfig, tasks *shutdown.Tasks) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")

	tasks.AddShutdownTask(
		func(ctx context.Context) error {
			return rdb.Close()
		},
	)

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}
