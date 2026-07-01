package common

import "fmt"

// GetUserTokenRedisKey 用户会话 token 的 Redis key。
// 同一 userID 只存一个 token，实现单设备语义 + 服务端可主动吊销。
func GetUserTokenRedisKey(userID int64) string {
	return fmt.Sprintf("token:%d", userID)
}
