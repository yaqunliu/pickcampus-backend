package types

// RouteQuery /api/v1/route 入参（query 解析后）。坐标为 GCJ-02（与高德一致）。
type RouteQuery struct {
	OLng  float64 // 出发（家乡）经度
	OLat  float64 // 出发纬度
	DLng  float64 // 目的（院校）经度
	DLat  float64 // 目的纬度
	OCity string  // 出发城市名（高德跨城公交需要）
	DCity string  // 目的城市名
}

// RouteData /api/v1/route 返回 data：驾车 + 跨城公交（含高铁）。
// 指针字段：高德算不到某项时为 nil（omitempty 不下发），前端据此决定是否显示。
type RouteData struct {
	DrivingMin     *int  `json:"driving_min,omitempty"`      // 驾车时长(分钟)
	DrivingKm      *int  `json:"driving_km,omitempty"`       // 驾车里程(公里)
	TransitMin     *int  `json:"transit_min,omitempty"`      // 公交/高铁时长(分钟)
	TransitKm      *int  `json:"transit_km,omitempty"`       // 公交/高铁里程(公里)
	TransitHasRail *bool `json:"transit_has_rail,omitempty"` // 该方案是否含高铁/动车段
}
