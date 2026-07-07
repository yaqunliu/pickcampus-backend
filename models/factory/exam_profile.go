package factory

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"pickcampus-backend/models"
	"pickcampus-backend/models/repo"
)

type examProfileCrudImpl struct {
	db *gorm.DB
}

// ExamProfileRepo 构造测档档案 repo 实现，db 可为主库或事务。
func ExamProfileRepo(db *gorm.DB) repo.ExamProfileRepo {
	return &examProfileCrudImpl{db: db}
}

func (r *examProfileCrudImpl) Get(userID int64) (*models.ExamProfile, error) {
	var p models.ExamProfile
	err := r.db.Where("user_id = ?", userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrExamProfileNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *examProfileCrudImpl) Upsert(profile *models.ExamProfile) error {
	// 命中 user_id 唯一键则更新除主键/创建时间外的业务列，实现「一用户一行」整体覆盖。
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"score_mode", "province", "subject", "electives", "score", "rank", "update_time",
		}),
	}).Create(profile).Error
}
