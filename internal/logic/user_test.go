package logic

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/config"
	"pickcampus-backend/internal/types"
)

// setupJWTConfig 为不依赖真实配置文件的单测注入 JWT 配置。
func setupJWTConfig() {
	config.G.Conf.JWTTokenSecret = "test-secret-for-unit-test"
	config.G.Conf.JWTTokenExpires = 3600
}

// TestBcryptHashAndCompare 验证密码哈希与比对：明文不入库、正确密码通过、错误密码拒绝。
func TestBcryptHashAndCompare(t *testing.T) {
	password := "S3cur3Pass!"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	// 哈希与明文不同（未明文入库）
	assert.NotEqual(t, password, string(hash))

	// 正确密码通过
	assert.NoError(t, bcrypt.CompareHashAndPassword(hash, []byte(password)))
	// 错误密码拒绝
	assert.Error(t, bcrypt.CompareHashAndPassword(hash, []byte("wrong-password")))
}

// TestGenerateAndParseToken 验证 JWT 签发与解析往返一致。
func TestGenerateAndParseToken(t *testing.T) {
	setupJWTConfig()

	token, expiresIn, err := common.GenerateToken(42, "alice@example.com", common.UserRoleUser, common.UserStatusActive)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, int64(3600), expiresIn)

	claims, err := common.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), claims.UserID)
	assert.Equal(t, "alice@example.com", claims.Email)
	assert.Equal(t, common.UserRoleUser, claims.Role)
	assert.Equal(t, common.UserStatusActive, claims.Status)
}

// TestParseTokenRejectsTampered 验证篡改后的 token 被拒。
func TestParseTokenRejectsTampered(t *testing.T) {
	setupJWTConfig()

	token, _, err := common.GenerateToken(1, "bob@example.com", common.UserRoleUser, common.UserStatusActive)
	assert.NoError(t, err)

	// 篡改末尾一位破坏签名
	tampered := token[:len(token)-1] + "x"
	_, err = common.ParseToken(tampered)
	assert.Error(t, err)
}

// TestParseTokenRejectsExpired 验证过期 token 被拒。
// 直接用相同密钥手工签一个 ExpiresAt 在过去的 token，绕过 GenerateToken 的兜底逻辑。
func TestParseTokenRejectsExpired(t *testing.T) {
	setupJWTConfig()

	claims := &common.Claims{
		UserID: 1,
		Email:  "carol@example.com",
		Role:   common.UserRoleUser,
		Status: common.UserStatusActive,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.G.Conf.JWTTokenSecret))
	assert.NoError(t, err)

	_, err = common.ParseToken(signed)
	assert.Error(t, err)
}

// TestRegisterRequestValidation 验证注册请求的 binding 校验规则。
func TestRegisterRequestValidation(t *testing.T) {
	// gin 用 "binding" tag 触发校验，独立 validator 默认读 "validate" tag，
	// 这里显式改为读 "binding"，以真实反映 gin 的校验行为。
	v := validator.New()
	v.SetTagName("binding")

	cases := []struct {
		name    string
		req     types.RegisterRequest
		wantErr bool
	}{
		{"合法请求", types.RegisterRequest{Email: "u@example.com", Password: "12345678", Username: "u"}, false},
		{"username 可选为空", types.RegisterRequest{Email: "u@example.com", Password: "12345678"}, false},
		{"邮箱格式不校验（可随便填）", types.RegisterRequest{Email: "not-an-email", Password: "12345678"}, false},
		{"邮箱缺失", types.RegisterRequest{Password: "12345678"}, true},
		{"密码过短", types.RegisterRequest{Email: "u@example.com", Password: "123"}, true},
		{"密码缺失", types.RegisterRequest{Email: "u@example.com"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Struct(tc.req)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
