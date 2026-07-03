package common

// 错误码集中定义，按业务域分段。
const (
	// HTTP 通用段
	ErrCodeSuccess         = 0
	ErrCodeInvalidInput    = 400
	ErrCodeUnauthorized    = 401
	ErrCodeForbidden       = 403
	ErrCodeNotFound        = 404
	ErrCodeTooManyRequests = 429
	ErrCodeInternalError   = 500

	// 用户段（1000-1099）
	ErrCodeEmailExists        = 1001 // 邮箱已存在
	ErrCodeUserNotFound       = 1003 // 用户不存在
	ErrCodeInvalidCredentials = 1004 // 凭证无效（邮箱或密码错误，防枚举）

	// 候选段（1100-1199）
	ErrCodeCandidateNotFound = 1101 // 候选不存在（删除/加专业时该校未在候选）
	ErrCodeCandidateExists   = 1102 // 候选已存在（幂等场景一般不报错，预留）

	// 数据库段（5000-5099）
	ErrCodeDatabaseError = 5000

	// 加密/token 段（5300-5399）
	ErrCodeCryptoError = 5300
	ErrCodeTokenError  = 5301
)
