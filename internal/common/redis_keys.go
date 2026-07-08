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

// AdmissionCacheTTL 录取/专业数据缓存有效期：24 小时。
// 数据仅在重导入时变化，长缓存显著降低 DB 压力；重导入后最多 24h 内自然刷新。
const AdmissionCacheTTL = 24 * time.Hour

// GetAdmissionRedisKey 按省录取列表的 Redis key。majorLevel 区分院校级/专业级。
func GetAdmissionRedisKey(province string, majorLevel bool) string {
	level := "college"
	if majorLevel {
		level = "major"
	}
	return fmt.Sprintf("admission:%s:%s", level, province)
}

// GetUniversityMajorsRedisKey 按校开设专业名清单的 Redis key。
func GetUniversityMajorsRedisKey(schoolID string) string {
	return fmt.Sprintf("univ-majors:%s", schoolID)
}
