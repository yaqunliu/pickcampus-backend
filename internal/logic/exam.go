package logic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"pickcampus-backend/internal/bootstrap"
	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/types"
	"pickcampus-backend/models"
	"pickcampus-backend/models/factory"
	"pickcampus-backend/models/repo"
)

// ExamLogic 测档档案业务逻辑。
type ExamLogic struct {
	Ctx context.Context
	DB  *gorm.DB
}

// NewExamLogic 构造，从单例拿 db。
func NewExamLogic(ctx context.Context) *ExamLogic {
	return &ExamLogic{Ctx: ctx, DB: bootstrap.Cli(ctx)}
}

// Get 取用户测档档案；无档案返回 (nil, nil)，交由 handler 输出空。
func (l *ExamLogic) Get(userID int64) (*types.ExamProfileDTO, error) {
	p, err := factory.ExamProfileRepo(l.DB).Get(userID)
	if errors.Is(err, repo.ErrExamProfileNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, NewBizError(common.ErrCodeDatabaseError, "数据库错误")
	}
	return toExamDTO(p), nil
}

// Save 整体保存(upsert)用户测档档案。
func (l *ExamLogic) Save(userID int64, dto types.ExamProfileDTO) error {
	electives := models.MajorList(dto.Electives)
	if electives == nil {
		electives = models.MajorList{}
	}
	profile := &models.ExamProfile{
		UserID:    userID,
		ScoreMode: dto.ScoreMode,
		Province:  dto.Province,
		Subject:   dto.Subject,
		Electives: electives,
		Score:     dto.Score,
		Rank:      dto.Rank,
	}
	if err := factory.ExamProfileRepo(l.DB).Upsert(profile); err != nil {
		return NewBizError(common.ErrCodeDatabaseError, "保存测档失败")
	}
	return nil
}

// toExamDTO 模型→对外 DTO。
func toExamDTO(p *models.ExamProfile) *types.ExamProfileDTO {
	return &types.ExamProfileDTO{
		ScoreMode: p.ScoreMode,
		Province:  p.Province,
		Subject:   p.Subject,
		Electives: []string(p.Electives),
		Score:     p.Score,
		Rank:      p.Rank,
	}
}
