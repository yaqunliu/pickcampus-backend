package factory

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"pickcampus-backend/models"
	"pickcampus-backend/models/repo"
)

type candidateCrudImpl struct {
	db *gorm.DB
}

// CandidateRepo 构造候选 repo 实现，db 可为主库或事务。
func CandidateRepo(db *gorm.DB) repo.CandidateRepo {
	return &candidateCrudImpl{db: db}
}

func (r *candidateCrudImpl) ListByUser(userID int64) ([]*models.Candidate, error) {
	var list []*models.Candidate
	err := r.db.Where("user_id = ?", userID).Order("create_time asc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *candidateCrudImpl) Get(userID int64, schoolID string) (*models.Candidate, error) {
	var c models.Candidate
	err := r.db.Where("user_id = ? AND school_id = ?", userID, schoolID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrCandidateNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *candidateCrudImpl) EnsureSchool(userID int64, schoolID string) (*models.Candidate, error) {
	c := &models.Candidate{
		UserID:   userID,
		SchoolID: schoolID,
		Majors:   models.MajorList{},
	}
	// 命中 (user_id, school_id) 唯一键则什么都不做，实现幂等新增。
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "school_id"}},
		DoNothing: true,
	}).Create(c).Error
	if err != nil {
		return nil, err
	}
	// OnConflict DoNothing 时 c 可能未回填最新行，统一再查一次返回准确数据。
	return r.Get(userID, schoolID)
}

func (r *candidateCrudImpl) Delete(userID int64, schoolID string) error {
	// 目标不存在时 Delete 影响 0 行、返回 nil，天然幂等。
	return r.db.Where("user_id = ? AND school_id = ?", userID, schoolID).
		Delete(&models.Candidate{}).Error
}

func (r *candidateCrudImpl) UpdateMajors(userID int64, schoolID string, majors models.MajorList) error {
	return r.db.Model(&models.Candidate{}).
		Where("user_id = ? AND school_id = ?", userID, schoolID).
		Update("majors", majors).Error
}
