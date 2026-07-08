package repo

import "pickcampus-backend/models"

// AdmissionRepo 录取数据访问接口（院校级 + 专业级同表，major 是否为空区分）。
type AdmissionRepo interface {
	// DeleteAll 清空全表（重导入前）。
	DeleteAll() error
	// BulkInsert 批量插入。
	BulkInsert(records []*models.Admission) error
	// ListByProvince 取某省记录；majorLevel=false 取院校级(major IS NULL)，true 取专业级(major IS NOT NULL)。
	ListByProvince(province string, majorLevel bool) ([]*models.Admission, error)
	// DistinctMajorsBySchool 取某校去重后的开设专业名（major IS NOT NULL）。
	DistinctMajorsBySchool(universityID string) ([]string, error)
	// CountByLevel 统计院校级/专业级行数（导入后核对用）。
	CountByLevel(majorLevel bool) (int64, error)
}
