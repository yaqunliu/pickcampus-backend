package models

// User 用户表（tbl_user）。邮箱为唯一登录标识，密码存 bcrypt 哈希。
type User struct {
	Meta
	Email        string `json:"email" gorm:"column:email;type:varchar(128);not null;uniqueIndex:idx_email;comment:邮箱(唯一登录标识)"`
	PasswordHash string `json:"-" gorm:"column:password_hash;type:varchar(128);not null;comment:bcrypt密码哈希"` // json:"-" 绝不外泄
	Username     string `json:"username" gorm:"column:username;type:varchar(64);default:'';comment:展示名(可选,空则用邮箱)"`
	Role         string `json:"role" gorm:"column:role;type:varchar(32);default:user;not null;index:idx_role;comment:角色(user/admin)"`
	Status       string `json:"status" gorm:"column:status;type:varchar(32);default:active;not null;index:idx_status;comment:状态(active/disabled)"`
}
