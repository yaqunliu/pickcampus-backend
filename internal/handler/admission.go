package handler

import (
	"github.com/gin-gonic/gin"

	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/logic"
)

// AdmissionHandler 录取/专业数据 HTTP 处理器（公开只读）。
type AdmissionHandler struct{}

// NewAdmissionHandler 构造。
func NewAdmissionHandler() *AdmissionHandler {
	return &AdmissionHandler{}
}

// College 取某省院校级录取列表。
func (h *AdmissionHandler) College(c *gin.Context) {
	h.listByProvince(c, false)
}

// Major 取某省专业级录取列表。
func (h *AdmissionHandler) Major(c *gin.Context) {
	h.listByProvince(c, true)
}

func (h *AdmissionHandler) listByProvince(c *gin.Context, majorLevel bool) {
	province := c.Query("province")
	if province == "" {
		Error(c, common.ErrCodeInvalidInput, "缺少 province 参数")
		return
	}
	data, err := logic.NewAdmissionLogic(c.Request.Context()).ListByProvince(province, majorLevel)
	if err != nil {
		HandleError(c, err)
		return
	}
	setPublicCache(c)
	Success(c, "成功", data)
}

// UniversityMajors 取某校去重后的开设专业名清单。
func (h *AdmissionHandler) UniversityMajors(c *gin.Context) {
	schoolID := c.Query("school_id")
	if schoolID == "" {
		Error(c, common.ErrCodeInvalidInput, "缺少 school_id 参数")
		return
	}
	data, err := logic.NewAdmissionLogic(c.Request.Context()).MajorsBySchool(schoolID)
	if err != nil {
		HandleError(c, err)
		return
	}
	setPublicCache(c)
	Success(c, "成功", data)
}

// setPublicCache 录取数据只读、变动稀少,允许浏览器/CDN 缓存 1 天,进一步降后端压力。
func setPublicCache(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=86400")
}
