package common

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"pickcampus-backend/internal/config"
)

// Claims 自定义 JWT 载荷。
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
	jwt.RegisteredClaims
}

// getJWTSecret 从全局配置取签名密钥，空则 panic（配置缺失属致命错误）。
func getJWTSecret() []byte {
	secret := config.G.Conf.JWTTokenSecret
	if secret == "" {
		panic("jwt_token_secret 未配置")
	}
	return []byte(secret)
}

// tokenExpiration 从配置读取过期时长（秒），未配置则默认 24h。
func tokenExpiration() time.Duration {
	expires := config.G.Conf.JWTTokenExpires
	if expires <= 0 {
		expires = 86400
	}
	return time.Duration(expires) * time.Second
}

// GenerateToken 签发 JWT，返回 token 字符串与过期秒数。
func GenerateToken(userID int64, email, role, status string) (string, int64, error) {
	exp := tokenExpiration()
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Status: status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", 0, err
	}
	return signed, int64(exp.Seconds()), nil
}

// ParseToken 解析并校验 JWT：强制 HMAC 签名方法，校验签名与有效期。
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// 只接受 HMAC，防止算法混淆攻击
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("非预期的签名方法")
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token 无效")
	}
	// 兜底再显式比对过期时间
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}
	return claims, nil
}
