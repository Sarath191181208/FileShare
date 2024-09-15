package cache

import (
	"time"

	"github.com/go-redis/redis"
)

type Cache struct{
  Client     *redis.Client
}

func (c *Cache) Set(key string, value interface{}, expirationTime time.Duration) error {
  return c.Client.Set(key, value, expirationTime).Err()
}

func (c *Cache) Get(key string) (string, error) {
  return c.Client.Get(key).Result()
}


