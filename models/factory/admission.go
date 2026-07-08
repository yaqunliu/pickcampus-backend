package factory

import (
	"gorm.io/gorm"

	"pickcampus-backend/models"
	"pickcampus-backend/models/repo"
)

type admissionCrudImpl struct {
	db *gorm.DB
}

// AdmissionRepo 构造录取数据 repo 实现，db 可为主库或事务。
func AdmissionRepo(db *gorm.DB) repo.AdmissionRepo {
	return &admissionCrudImpl{db: db}
}

func (r *admissionCrudImpl) DeleteAll() error {
	// GORM 默认禁止无条件全表删除，用恒真条件放行。
	return r.db.Where("1 = 1").Delete(&models.Admission{}).Error
}

func (r *admissionCrudImpl) BulkInsert(records []*models.Admission) error {
	if len(records) == 0 {
		return nil
	}
	return r.db.CreateInBatches(records, 500).Error
}

func (r *admissionCrudImpl) ListByProvince(province string, majorLevel bool) ([]*models.Admission, error) {
	var list []*models.Admission
	q := r.db.Where("province = ?", province)
	if majorLevel {
		q = q.Where("major IS NOT NULL")
	} else {
		q = q.Where("major IS NULL")
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *admissionCrudImpl) DistinctMajorsBySchool(universityID string) ([]string, error) {
	var majors []string
	err := r.db.Model(&models.Admission{}).
		Where("university_id = ? AND major IS NOT NULL", universityID).
		Distinct().
		Pluck("major", &majors).Error
	if err != nil {
		return nil, err
	}
	return majors, nil
}

func (r *admissionCrudImpl) CountByLevel(majorLevel bool) (int64, error) {
	var n int64
	q := r.db.Model(&models.Admission{})
	if majorLevel {
		q = q.Where("major IS NOT NULL")
	} else {
		q = q.Where("major IS NULL")
	}
	return n, q.Count(&n).Error
}
