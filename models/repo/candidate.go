package repo

import (
	"errors"

	"pickcampus-backend/models"
)

// ErrCandidateNotFound 候选不存在的哨兵错误（供 logic 层判定）。
var ErrCandidateNotFound = errors.New("candidate not found")

// CandidateRepo 候选数据访问接口。
type CandidateRepo interface {
	// ListByUser 拉某用户全部候选（按创建时间升序）。
	ListByUser(userID int64) ([]*models.Candidate, error)
	// Get 取单条候选；未命中回 ErrCandidateNotFound。
	Get(userID int64, schoolID string) (*models.Candidate, error)
	// EnsureSchool 幂等新增院校候选行：已存在则不动，返回该行。
	EnsureSchool(userID int64, schoolID string) (*models.Candidate, error)
	// Delete 删除院校候选（连带其专业，天然级联）；不存在也返回 nil（幂等）。
	Delete(userID int64, schoolID string) error
	// UpdateMajors 更新某院校候选的专业数组。
	UpdateMajors(userID int64, schoolID string, majors models.MajorList) error
}
