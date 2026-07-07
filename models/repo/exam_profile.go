package repo

import (
	"errors"

	"pickcampus-backend/models"
)

// ErrExamProfileNotFound 测档档案不存在的哨兵错误（供 logic 层判定「首次、无档案」）。
var ErrExamProfileNotFound = errors.New("exam profile not found")

// ExamProfileRepo 测档档案数据访问接口（一用户一行）。
type ExamProfileRepo interface {
	// Get 取某用户的测档档案；未命中回 ErrExamProfileNotFound。
	Get(userID int64) (*models.ExamProfile, error)
	// Upsert 整体保存（存在则更新，不存在则插入）。
	Upsert(profile *models.ExamProfile) error
}
