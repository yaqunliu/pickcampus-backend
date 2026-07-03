package handler

import (
	"github.com/gin-gonic/gin"

	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/logic"
	"pickcampus-backend/internal/types"
)

// CandidateHandler 候选相关 HTTP 处理器。
type CandidateHandler struct{}

// NewCandidateHandler 构造。
func NewCandidateHandler() *CandidateHandler {
	return &CandidateHandler{}
}

// List 拉当前用户全部候选。
func (h *CandidateHandler) List(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	items, err := logic.NewCandidateLogic(c.Request.Context()).List(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessWithTotal(c, "成功", items, int64(len(items)))
}

// AddSchool 加入候选院校（幂等）。
func (h *CandidateHandler) AddSchool(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	var req types.AddCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, common.ErrCodeInvalidInput, "请求参数不合法")
		return
	}
	if err := logic.NewCandidateLogic(c.Request.Context()).AddSchool(userID, req.SchoolID); err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "已加入候选", nil)
}

// RemoveSchool 移除候选院校（级联其专业）。
func (h *CandidateHandler) RemoveSchool(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	schoolID := c.Param("school_id")
	if schoolID == "" {
		Error(c, common.ErrCodeInvalidInput, "缺少院校ID")
		return
	}
	if err := logic.NewCandidateLogic(c.Request.Context()).RemoveSchool(userID, schoolID); err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "已移除候选", nil)
}

// AddMajor 加候选专业（校不存在则自动 upsert）。
func (h *CandidateHandler) AddMajor(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	schoolID := c.Param("school_id")
	if schoolID == "" {
		Error(c, common.ErrCodeInvalidInput, "缺少院校ID")
		return
	}
	var req types.CandidateMajorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, common.ErrCodeInvalidInput, "请求参数不合法")
		return
	}
	if err := logic.NewCandidateLogic(c.Request.Context()).AddMajor(userID, schoolID, req.Major); err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "已加入候选专业", nil)
}

// RemoveMajor 移除候选专业。
func (h *CandidateHandler) RemoveMajor(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	schoolID := c.Param("school_id")
	if schoolID == "" {
		Error(c, common.ErrCodeInvalidInput, "缺少院校ID")
		return
	}
	var req types.CandidateMajorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, common.ErrCodeInvalidInput, "请求参数不合法")
		return
	}
	if err := logic.NewCandidateLogic(c.Request.Context()).RemoveMajor(userID, schoolID, req.Major); err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "已移除候选专业", nil)
}
