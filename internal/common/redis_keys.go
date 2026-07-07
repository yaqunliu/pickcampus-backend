package common

import (
	"fmt"
	"time"
)

// GetUserTokenRedisKey 用户会话 token 的 Redis key。
// 同一 userID 只存一个 token，实现单设备语义 + 服务端可主动吊销。
func GetUserTokenRedisKey(userID int64) string {
	return fmt.Sprintf("token:%d", userID)
}

// RouteCacheTTL 路程查询结果缓存有效期：30 天。
// 高铁/驾车时长基本不随时间变化，长缓存即可，减少高德调用。
const RouteCacheTTL = 30 * 24 * time.Hour

// GetRouteCacheRedisKey 路程查询结果的 Redis key。
// 由出发/目的坐标（6 位小数）与城市名唯一确定同一次查询。
func GetRouteCacheRedisKey(oLng, oLat, dLng, dLat float64, oCity, dCity string) string {
	return fmt.Sprintf("route:%.6f,%.6f:%.6f,%.6f:%s:%s", oLng, oLat, dLng, dLat, oCity, dCity)
}
