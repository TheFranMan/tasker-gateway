package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"gateway/common"
)

type Interface interface {
	GetKey(key string) (string, error)
	SetKey(name string, value interface{}) error
}

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func New(config *common.Config) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: "",
			DB:       0,
		}),
		ttl: config.RedisKeyTtl,
	}
}

func (c *Cache) SetKey(name string, value interface{}) error {
	return c.client.Set(context.Background(), name, value, c.ttl).Err()
}

func (c *Cache) GetKey(key string) (string, error) {
	status, err := c.client.Get(context.Background(), key).Result()
	if nil != err {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", err
	}

	return status, nil
}
