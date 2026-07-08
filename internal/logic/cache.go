package logic

import (
	"context"
	"encoding/json"
	"time"

	"pickcampus-backend/internal/bootstrap"
)

// cacheGetJSON 从 Redis 读并反序列化为 T。未初始化/未命中/解析失败均返回 (zero, false),视为未命中。
func cacheGetJSON[T any](ctx context.Context, key string) (T, bool) {
	var zero T
	cli := bootstrap.GetCli()
	if cli == nil {
		return zero, false
	}
	raw, err := bootstrap.NewCRUD(ctx, cli).Get(key)
	if err != nil {
		return zero, false // redis.Nil(未命中)或其它错误,均降级到 DB
	}
	var v T
	if json.Unmarshal([]byte(raw), &v) != nil {
		return zero, false
	}
	return v, true
}

// cacheSetJSON 序列化并写 Redis(带 TTL)。失败仅忽略,不影响本次返回。
func cacheSetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) {
	cli := bootstrap.GetCli()
	if cli == nil {
		return
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return
	}
	_ = bootstrap.NewCRUD(ctx, cli).Set(key, raw, ttl)
}
