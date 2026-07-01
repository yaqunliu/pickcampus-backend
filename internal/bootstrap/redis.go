package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"pickcampus-backend/internal/config"
)

// _rdb Redis 单例。
var _rdb *redis.Client

// InitRedis 建立 Redis 连接并 Ping 校验连通性。
func InitRedis(cfg config.RedisConfig) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}
	_rdb = rdb
	return nil
}

// GetCli 返回 Redis 单例，nil 表示未初始化。
func GetCli() *redis.Client {
	return _rdb
}

// RedisCrud 精简的 Redis 读写封装。
type RedisCrud struct {
	ctx context.Context
	cli *redis.Client
}

// NewCRUD 构造 CRUD 封装。
func NewCRUD(ctx context.Context, cli *redis.Client) *RedisCrud {
	return &RedisCrud{ctx: ctx, cli: cli}
}

// Set 写入并设置过期时间。
func (c *RedisCrud) Set(key string, value interface{}, timeout time.Duration) error {
	return c.cli.Set(c.ctx, key, value, timeout).Err()
}

// Get 读取；未命中返回 redis.Nil 错误（由调用方用 errors.Is 判定）。
func (c *RedisCrud) Get(key string) (string, error) {
	return c.cli.Get(c.ctx, key).Result()
}

// Delete 删除 key。
func (c *RedisCrud) Delete(key string) error {
	return c.cli.Del(c.ctx, key).Err()
}
