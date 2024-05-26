package storage

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// Redis setter wrapper
func Set(c *redis.Client, key string, value interface{}, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(key, p, expiration).Err()
}

// Redis getter wrapper
func Get(c *redis.Client, key string, dest interface{}) error {
	p, err := c.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(p), &dest)
}

func Delete(c *redis.Client, key string) {
	c.Del(key)
}
