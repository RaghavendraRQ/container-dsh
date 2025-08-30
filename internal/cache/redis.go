package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/raghavendrarq/container-dsh/internal/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb          *redis.Client
	containerKey string
}

func NewCache(containerKey string) (*Cache, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read redis config: %w", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}
	return &Cache{rdb: rdb, containerKey: containerKey}, nil
}

func (c *Cache) GetContainerList(ctx context.Context) ([]string, error) {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("can't ping to redis in getting")
	}
	values, err := c.rdb.LRange(ctx, c.containerKey, 0, -1).Result()
	if err == redis.Nil || len(values) == 0 {
		return nil, fmt.Errorf("can't get the containerids")
	} else if err != nil {
		return nil, fmt.Errorf("network issue or redis down")
	}
	return values, nil
}

func (c *Cache) Ping(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("can't ping to redis")
	}
	return nil
}

func (c *Cache) SetContainerList(ctx context.Context, data []string, ttl time.Duration) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("can't ping to redis")
	}
	// Batching all the commands and sending them
	pipe := c.rdb.Pipeline()
	pipe.Del(ctx, c.containerKey)
	for _, data := range data {
		pipe.RPush(ctx, c.containerKey, data)
	}
	pipe.Expire(ctx, c.containerKey, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("can't get the containerids")
	}
	return nil
}
