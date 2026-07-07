package handler

import (
	"github.com/gin-gonic/gin"

	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/logic"
	"pickcampus-backend/internal/types"
)

// ExamHandler 测档档案 HTTP 处理器。
type ExamHandler struct{}

// NewExamHandler 构造。
func NewExamHandler() *ExamHandler {
	return &ExamHandler{}
}

// Get 取当前用户测档档案；无档案时 data 为空。
func (h *ExamHandler) Get(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	dto, err := logic.NewExamLogic(c.Request.Context()).Get(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "成功", dto) // dto 为 nil 时 data 省略,前端视为无档案
}

// Save 整体保存(upsert)当前用户测档档案。
func (h *ExamHandler) Save(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	var req types.ExamProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, common.ErrCodeInvalidInput, "请求参数不合法")
		return
	}
	if err := logic.NewExamLogic(c.Request.Context()).Save(userID, req); err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "已保存", nil)
}
