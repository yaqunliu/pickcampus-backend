// Package logic 承载业务编排（事务、bcrypt、token、Redis 会话）。
package logic

import "fmt"

// BizError 业务错误，携带错误码，由 handler 的 HandleError 翻译成响应。
type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return fmt.Sprintf("biz error %d: %s", e.Code, e.Message)
}

// NewBizError 构造业务错误。
func NewBizError(code int, message string) *BizError {
	return &BizError{Code: code, Message: message}
}
