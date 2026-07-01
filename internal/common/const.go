package common

// 用户角色与状态字面量常量（不使用 TS 风格枚举，直接常量）。
const (
	UserRoleUser  = "user"
	UserRoleAdmin = "admin"

	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)
