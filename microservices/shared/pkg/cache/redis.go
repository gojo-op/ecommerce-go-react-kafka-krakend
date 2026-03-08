package cache

import (
    "context"
    "time"
    cfgpkg "github.com/your-org/microservices/shared/config"
)

type RedisClient struct{}

func NewRedisClient(cfg *cfgpkg.Config) (*RedisClient, error) { return &RedisClient{}, nil }
func (c *RedisClient) Close() error { return nil }
func (c *RedisClient) Delete(ctx context.Context, key string) error { return nil }
func (c *RedisClient) SetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) error { return nil }
func (c *RedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error { return fmtErr() }

func fmtErr() error { return nil }