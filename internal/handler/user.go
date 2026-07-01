package handler

import (
	"github.com/gin-gonic/gin"

	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/logic"
	"pickcampus-backend/internal/types"
)

// UserHandler 用户相关 HTTP 处理器。
type UserHandler struct{}

// NewUserHandler 构造。
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Register 注册。
func (h *UserHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, common.ErrCodeInvalidInput, "请求参数不合法")
		return
	}
	resp, err := logic.NewUserLogic(c.Request.Context()).Register(&req)
	if err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "注册成功", resp)
}

// Login 登录。
func (h *UserHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, common.ErrCodeInvalidInput, "请求参数不合法")
		return
	}
	resp, err := logic.NewUserLogic(c.Request.Context()).Login(&req)
	if err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "登录成功", resp)
}

// Logout 登出（吊销 Redis 会话）。
func (h *UserHandler) Logout(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	if err := logic.NewUserLogic(c.Request.Context()).Logout(userID); err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "登出成功", nil)
}

// GetUserInfo 拿当前用户。
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, ok := common.GetUserID(c)
	if !ok {
		Error(c, common.ErrCodeUnauthorized, "未登录")
		return
	}
	resp, err := logic.NewUserLogic(c.Request.Context()).GetUserInfo(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "成功", resp)
}
