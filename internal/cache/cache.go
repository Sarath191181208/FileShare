package cache

import (
	"errors"
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


type MockCache struct{}


func (m *MockCache) Set(key string, value interface{}, time time.Duration) error {
  return nil
}

func (m *MockCache) Get(key string) (string, error) {
  return "", errors.New("not found")
}

func (m *MockCache) Delete(key string) error {
  return nil
}
