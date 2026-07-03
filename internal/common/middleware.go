package common

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"pickcampus-backend/internal/bootstrap"
)

// abort401 统一以 401 + {code,message,data} 中断请求。
func abort401(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    ErrCodeUnauthorized,
		"message": message,
	})
}

// AuthMiddleware JWT + Redis 会话强校验（fail-secure）。
// 流程：解析 Authorization: Bearer <token> → 验签验期 → 从 Redis 读 token:{uid}
// 并与请求 token 逐字比对；不一致 / 不存在 / Redis 不可用一律 401。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abort401(c, "缺少 Authorization 头")
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			abort401(c, "Authorization 头格式应为 Bearer <token>")
			return
		}
		token := parts[1]

		claims, err := ParseToken(token)
		if err != nil {
			abort401(c, "token 无效或已过期")
			return
		}

		// 强制 Redis 校验（fail-secure）
		cli := bootstrap.GetCli()
		if cli == nil {
			abort401(c, "会话服务不可用")
			return
		}
		storedToken, err := bootstrap.NewCRUD(c.Request.Context(), cli).Get(GetUserTokenRedisKey(claims.UserID))
		if errors.Is(err, redis.Nil) {
			abort401(c, "会话已失效，请重新登录")
			return
		}
		if err != nil {
			abort401(c, "会话校验失败")
			return
		}
		if storedToken != token {
			abort401(c, "token 已被吊销")
			return
		}

		// 注入 context 供下游 handler 使用
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("status", claims.Status)
		c.Next()
	}
}

// GetUserID 从 context 取当前用户 ID。
func GetUserID(c *gin.Context) (int64, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}

// GetEmail 从 context 取当前用户邮箱。
func GetEmail(c *gin.Context) (string, bool) {
	v, exists := c.Get("email")
	if !exists {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// Cors 跨域中间件（前后端分离，允许携带凭证）。
// allowOriginsCSV 为逗号分隔的允许来源清单：
//   - 空:回显任意 Origin（本地开发,localhost 任意端口可连）
//   - 非空:仅回显命中清单的 Origin（生产收敛到前端域名,未命中则不发 Allow-Origin,浏览器拦截）
func Cors(allowOriginsCSV string) gin.HandlerFunc {
	allowed := make(map[string]struct{})
	for o := range strings.SplitSeq(allowOriginsCSV, ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowed[o] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			_, hit := allowed[origin]
			if len(allowed) == 0 || hit {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
