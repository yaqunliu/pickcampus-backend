package logic

import (
	"context"
	"errors"
	"slices"

	"gorm.io/gorm"

	"pickcampus-backend/internal/bootstrap"
	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/types"
	"pickcampus-backend/models"
	"pickcampus-backend/models/factory"
	"pickcampus-backend/models/repo"
)

// CandidateLogic 候选业务逻辑。
type CandidateLogic struct {
	Ctx context.Context
	DB  *gorm.DB
}

// NewCandidateLogic 构造，从单例拿 db。
func NewCandidateLogic(ctx context.Context) *CandidateLogic {
	return &CandidateLogic{Ctx: ctx, DB: bootstrap.Cli(ctx)}
}

// List 拉当前用户全部候选。
func (l *CandidateLogic) List(userID int64) ([]types.CandidateItem, error) {
	rows, err := factory.CandidateRepo(l.DB).ListByUser(userID)
	if err != nil {
		return nil, NewBizError(common.ErrCodeDatabaseError, "数据库错误")
	}
	items := make([]types.CandidateItem, 0, len(rows))
	for _, c := range rows {
		items = append(items, toCandidateItem(c))
	}
	return items, nil
}

// AddSchool 加入候选院校（幂等）。
func (l *CandidateLogic) AddSchool(userID int64, schoolID string) error {
	if _, err := factory.CandidateRepo(l.DB).EnsureSchool(userID, schoolID); err != nil {
		return NewBizError(common.ErrCodeDatabaseError, "加入候选失败")
	}
	return nil
}

// RemoveSchool 移除候选院校（连带其专业，天然级联）。
func (l *CandidateLogic) RemoveSchool(userID int64, schoolID string) error {
	if err := factory.CandidateRepo(l.DB).Delete(userID, schoolID); err != nil {
		return NewBizError(common.ErrCodeDatabaseError, "移除候选失败")
	}
	return nil
}

// AddMajor 加候选专业：校不存在则自动建院校行；majors 去重追加。
// 事务保护「读-改-写」，避免并发下丢更新。
func (l *CandidateLogic) AddMajor(userID int64, schoolID, major string) error {
	err := l.DB.Transaction(func(tx *gorm.DB) error {
		r := factory.CandidateRepo(tx)
		c, err := r.EnsureSchool(userID, schoolID)
		if err != nil {
			return err
		}
		if slices.Contains(c.Majors, major) {
			return nil // 已存在，幂等
		}
		next := append(models.MajorList{}, c.Majors...)
		next = append(next, major)
		return r.UpdateMajors(userID, schoolID, next)
	})
	if err != nil {
		return NewBizError(common.ErrCodeDatabaseError, "加候选专业失败")
	}
	return nil
}

// RemoveMajor 移除候选专业；该校不在候选则幂等返回成功。
func (l *CandidateLogic) RemoveMajor(userID int64, schoolID, major string) error {
	err := l.DB.Transaction(func(tx *gorm.DB) error {
		r := factory.CandidateRepo(tx)
		c, err := r.Get(userID, schoolID)
		if err != nil {
			if errors.Is(err, repo.ErrCandidateNotFound) {
				return nil // 校未在候选，无专业可删，幂等
			}
			return err
		}
		next := make(models.MajorList, 0, len(c.Majors))
		for _, m := range c.Majors {
			if m != major {
				next = append(next, m)
			}
		}
		if len(next) == len(c.Majors) {
			return nil // 无变化
		}
		return r.UpdateMajors(userID, schoolID, next)
	})
	if err != nil {
		return NewBizError(common.ErrCodeDatabaseError, "移除候选专业失败")
	}
	return nil
}

// toCandidateItem 领域模型转对外 DTO。majors 保证非 nil（前端直接消费）。
func toCandidateItem(c *models.Candidate) types.CandidateItem {
	majors := c.Majors
	if majors == nil {
		majors = models.MajorList{}
	}
	return types.CandidateItem{
		SchoolID:   c.SchoolID,
		Majors:     majors,
		CreateTime: c.CreateTime,
	}
}
