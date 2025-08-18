package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCfg struct {
	Addr     string
	Password string
}

func NewRedisConfig() (*redisCfg, error) {
	rdb := &redisCfg{}
	Addr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		rdb.Addr = "localhost:6379"
	} else {
		rdb.Addr = Addr
	}

	Password, ok := os.LookupEnv("REDIS_PASSWD")
	if !ok {
		rdb.Password = ""
	} else {
		rdb.Password = Password
	}

	return rdb, nil

}

type Cache struct {
	rdb          *redis.Client
	containerKey string
}

func NewCache(containerKey string) *Cache {
	redisClient, err := NewRedisConfig()
	if err != nil {
		panic("Can't Create Redis client")
	}
	rdbQ := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: redisClient.Password,
		DB:       0,
	})
	return &Cache{
		rdb:          rdbQ,
		containerKey: containerKey,
	}
}

func (c *Cache) GetContainerList(ctx context.Context) ([]string, error) {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("can't ping to redis in getting")
	}
	values, err := c.rdb.LRange(ctx, c.containerKey, 0, -1).Result()
	if err != nil || len(values) <= 0 {
		return nil, fmt.Errorf("can't get the containerids")
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
	// For Setting TTL on list Used Pipe
	pipe := c.rdb.Pipeline()
	pipe.RPush(ctx, c.containerKey, data)
	pipe.Expire(ctx, c.containerKey, ttl)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return fmt.Errorf("can't get the containerids")
	}
	return nil
}
