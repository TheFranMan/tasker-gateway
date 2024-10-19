package cache

import (
	"context"
	"errors"
	"time"

	"github.com/TheFranMan/tasker-common/types"
	"github.com/redis/go-redis/v9"

	"gateway/common"
)

type Interface interface {
	GetKey(key string) (*types.RequestStatusString, error)
	SetKey(key string, value types.RequestStatusString) error
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

func (c *Cache) SetKey(key string, value types.RequestStatusString) error {
	return c.client.Set(context.Background(), key, string(value), c.ttl).Err()
}

func (c *Cache) GetKey(key string) (*types.RequestStatusString, error) {
	value, err := c.client.Get(context.Background(), key).Result()
	if nil != err {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	status := types.RequestStatusString(value)
	return &status, nil
}
