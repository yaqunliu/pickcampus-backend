// Package factory 是 repo 接口的 GORM 实现。工厂函数注入 *gorm.DB（可传事务 tx）。
package factory

import (
	"errors"

	"gorm.io/gorm"

	"pickcampus-backend/models"
	"pickcampus-backend/models/repo"
)

type userCrudImpl struct {
	db *gorm.DB
}

// UserRepo 构造用户 repo 实现，db 可为主库或事务。
func UserRepo(db *gorm.DB) repo.UserRepo {
	return &userCrudImpl{db: db}
}

func (r *userCrudImpl) Create(user *models.User) error {
	err := r.db.Create(user).Error
	// 并发注册时 ExistsByEmail 可能都放行，最终由邮箱唯一索引兜底；
	// 统一翻译成 ErrDuplicateEmail 供 logic 层给出正确的「邮箱已注册」提示。
	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
		return repo.ErrDuplicateEmail
	}
	return err
}

func (r *userCrudImpl) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userCrudImpl) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userCrudImpl) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
