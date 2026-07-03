// Package types 集中定义请求/响应 DTO。
package types

// RegisterRequest 注册请求。username 可选。
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,max=128"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	Username string `json:"username" binding:"omitempty,max=64"`
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserInfo 对外用户信息（不含密码）。
type UserInfo struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

// LoginResponse 登录响应。
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresIn int64     `json:"expires_in"`
	UserInfo  *UserInfo `json:"user_info"`
}
