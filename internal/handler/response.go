// Package handler 是 HTTP 层：校验入参、调 logic、封装响应。
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/logic"
)

// Response 统一响应体：{code,message,data}（遵循根 CLAUDE.md 通用标准）。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseWithTotal 列表响应，带总数——列表永远包在对象里，不裸返数组。
type ResponseWithTotal struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Total   int64       `json:"total"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应。
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: common.ErrCodeSuccess, Message: message, Data: data})
}

// SuccessWithTotal 带总数的列表成功响应。
func SuccessWithTotal(c *gin.Context, message string, data interface{}, total int64) {
	c.JSON(http.StatusOK, ResponseWithTotal{Code: common.ErrCodeSuccess, Message: message, Total: total, Data: data})
}

// Error 错误响应（业务错误统一 HTTP 200 + body code）。
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message})
}

// HandleError 把 logic 层的 BizError 翻译成响应；非 BizError 归一为 500。
func HandleError(c *gin.Context, err error) {
	var bizErr *logic.BizError
	if errors.As(err, &bizErr) {
		Error(c, bizErr.Code, bizErr.Message)
		return
	}
	Error(c, common.ErrCodeInternalError, "服务器内部错误")
}
