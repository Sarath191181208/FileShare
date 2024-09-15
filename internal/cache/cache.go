package cache

import (
	"time"

	"github.com/go-redis/redis"
)

type Cache interface{
  Set(key string, value interface{}, expirationTime time.Duration) error 
  Get(key string) (string, error)
  Delete(key string) error
}

type RedisCache struct{
  Client     *redis.Client
}

func (c *RedisCache) Set(key string, value interface{}, expirationTime time.Duration) error {
  return c.Client.Set(key, value, expirationTime).Err()
}

func (c *RedisCache) Get(key string) (string, error) {
  return c.Client.Get(key).Result()
}

func (c *RedisCache) Delete(key string) error {
  return c.Client.Del(key).Err()
}


