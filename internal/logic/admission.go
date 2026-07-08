package logic

import (
	"context"

	"gorm.io/gorm"

	"pickcampus-backend/internal/bootstrap"
	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/types"
	"pickcampus-backend/models"
	"pickcampus-backend/models/factory"
)

// AdmissionLogic 录取/专业数据查询业务逻辑（读穿透 Redis 缓存）。
type AdmissionLogic struct {
	Ctx context.Context
	DB  *gorm.DB
}

// NewAdmissionLogic 构造。
func NewAdmissionLogic(ctx context.Context) *AdmissionLogic {
	return &AdmissionLogic{Ctx: ctx, DB: bootstrap.Cli(ctx)}
}

// ListByProvince 取某省录取列表。majorLevel=false 院校级,true 专业级。先读缓存,未命中查库回填。
func (l *AdmissionLogic) ListByProvince(province string, majorLevel bool) ([]types.AdmissionRecordDTO, error) {
	cacheKey := common.GetAdmissionRedisKey(province, majorLevel)
	if cached, ok := cacheGetJSON[[]types.AdmissionRecordDTO](l.Ctx, cacheKey); ok {
		return cached, nil
	}
	rows, err := factory.AdmissionRepo(l.DB).ListByProvince(province, majorLevel)
	if err != nil {
		return nil, NewBizError(common.ErrCodeDatabaseError, "数据库错误")
	}
	dtos := make([]types.AdmissionRecordDTO, 0, len(rows))
	for _, m := range rows {
		dtos = append(dtos, toAdmissionDTO(m))
	}
	cacheSetJSON(l.Ctx, cacheKey, dtos, common.AdmissionCacheTTL)
	return dtos, nil
}

// MajorsBySchool 取某校去重后的开设专业名。先读缓存,未命中查库回填。
func (l *AdmissionLogic) MajorsBySchool(schoolID string) ([]string, error) {
	cacheKey := common.GetUniversityMajorsRedisKey(schoolID)
	if cached, ok := cacheGetJSON[[]string](l.Ctx, cacheKey); ok {
		return cached, nil
	}
	majors, err := factory.AdmissionRepo(l.DB).DistinctMajorsBySchool(schoolID)
	if err != nil {
		return nil, NewBizError(common.ErrCodeDatabaseError, "数据库错误")
	}
	if majors == nil {
		majors = []string{}
	}
	cacheSetJSON(l.Ctx, cacheKey, majors, common.AdmissionCacheTTL)
	return majors, nil
}

// toAdmissionDTO 模型 → 对外 camelCase DTO。
func toAdmissionDTO(m *models.Admission) types.AdmissionRecordDTO {
	var subs []types.SubRatingDTO
	if len(m.SubRatings) > 0 {
		subs = make([]types.SubRatingDTO, 0, len(m.SubRatings))
		for _, s := range m.SubRatings {
			subs = append(subs, types.SubRatingDTO{Major: s.Major, Rating: s.Rating})
		}
	}
	return types.AdmissionRecordDTO{
		UniversityID: m.UniversityID,
		Province:     m.Province,
		Subject:      m.Subject,
		Year:         m.Year,
		MinRank:      m.MinRank,
		MinScore:     m.MinScore,
		Source:       m.Source,
		Major:        m.Major,
		MajorCode:    m.MajorCode,
		ElectiveReq:  m.ElectiveReq,
		Batch:        m.Batch,
		Tuition:      m.Tuition,
		Duration:     m.Duration,
		Rating:       m.Rating,
		RatingRank:   m.RatingRank,
		SubRatings:   subs,
	}
}
