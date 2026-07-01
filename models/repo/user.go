// Package repo 定义数据访问层接口。实现见 models/factory。
package repo

import (
	"errors"

	"pickcampus-backend/models"
)

// ErrUserNotFound 用户不存在的哨兵错误（供 logic 层判定）。
var ErrUserNotFound = errors.New("user not found")

// ErrDuplicateEmail 邮箱唯一键冲突的哨兵错误（并发注册竞态时由 DB 唯一索引兜底触发）。
var ErrDuplicateEmail = errors.New("duplicate email")

// UserRepo 用户数据访问接口。
type UserRepo interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	ExistsByEmail(email string) (bool, error)
}
