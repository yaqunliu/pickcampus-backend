package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/logic"
	"pickcampus-backend/internal/types"
)

// RouteHandler 路程查询 HTTP 处理器。
type RouteHandler struct{}

// NewRouteHandler 构造。
func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

// Query 查询家乡→院校路程（驾车 + 跨城公交）。公开接口，无需登录。
func (h *RouteHandler) Query(c *gin.Context) {
	q, err := parseRouteQuery(c)
	if err != nil {
		Error(c, common.ErrCodeInvalidInput, "请求参数不合法")
		return
	}
	data, err := logic.NewRouteLogic(c.Request.Context()).Query(q)
	if err != nil {
		HandleError(c, err)
		return
	}
	Success(c, "成功", data)
}

// parseRouteQuery 从 query 解析经纬度与城市名，任一缺失/非法即报错。
func parseRouteQuery(c *gin.Context) (types.RouteQuery, error) {
	olng, e1 := strconv.ParseFloat(c.Query("olng"), 64)
	olat, e2 := strconv.ParseFloat(c.Query("olat"), 64)
	dlng, e3 := strconv.ParseFloat(c.Query("dlng"), 64)
	dlat, e4 := strconv.ParseFloat(c.Query("dlat"), 64)
	ocity := c.Query("ocity")
	dcity := c.Query("dcity")
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || ocity == "" || dcity == "" {
		return types.RouteQuery{}, fmt.Errorf("invalid route params")
	}
	return types.RouteQuery{
		OLng: olng, OLat: olat, DLng: dlng, DLat: dlat, OCity: ocity, DCity: dcity,
	}, nil
}
