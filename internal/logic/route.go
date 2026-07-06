package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/config"
	"pickcampus-backend/internal/types"
)

// 高德 Web 服务路径规划 v3 端点。
const (
	amapDrivingURL = "https://restapi.amap.com/v3/direction/driving"
	amapTransitURL = "https://restapi.amap.com/v3/direction/transit/integrated"
)

// RouteLogic 路程查询业务逻辑（调高德，无 DB）。
type RouteLogic struct {
	Ctx context.Context
}

// NewRouteLogic 构造。
func NewRouteLogic(ctx context.Context) *RouteLogic {
	return &RouteLogic{Ctx: ctx}
}

// Query 查询家乡→院校路程：驾车 + 跨城公交（含高铁）。
// 两项各自独立，任一成功即返回；均失败才报上游错误（前端据此降级到直线估算）。
func (l *RouteLogic) Query(q types.RouteQuery) (*types.RouteData, error) {
	key := config.G.Conf.AmapKey
	if key == "" {
		return nil, NewBizError(common.ErrCodeRouteNotConfigured, "路程服务未配置")
	}

	client := &http.Client{Timeout: 6 * time.Second}
	origin := fmt.Sprintf("%.6f,%.6f", q.OLng, q.OLat)
	dest := fmt.Sprintf("%.6f,%.6f", q.DLng, q.DLat)

	data := &types.RouteData{}
	okAny := false
	if km, min, err := fetchDriving(client, key, origin, dest); err == nil {
		data.DrivingKm, data.DrivingMin = km, min
		okAny = true
	}
	if km, min, rail, err := fetchTransit(client, key, origin, dest, q.OCity, q.DCity); err == nil {
		data.TransitKm, data.TransitMin, data.TransitHasRail = km, min, rail
		okAny = true
	}
	if !okAny {
		return nil, NewBizError(common.ErrCodeRouteUpstream, "路程查询失败")
	}
	return data, nil
}

// fetchDriving 高德驾车路径规划：取首条方案的里程/时长。
func fetchDriving(client *http.Client, key, origin, dest string) (km, min *int, err error) {
	u := fmt.Sprintf("%s?origin=%s&destination=%s&extensions=base&key=%s",
		amapDrivingURL, url.QueryEscape(origin), url.QueryEscape(dest), key)
	var body struct {
		Status string `json:"status"`
		Route  struct {
			Paths []struct {
				Distance string `json:"distance"`
				Duration string `json:"duration"`
			} `json:"paths"`
		} `json:"route"`
	}
	if err = getJSON(client, u, &body); err != nil {
		return nil, nil, err
	}
	if body.Status != "1" || len(body.Route.Paths) == 0 {
		return nil, nil, fmt.Errorf("amap driving: no path")
	}
	p := body.Route.Paths[0]
	return kmFromMeters(p.Distance), minFromSeconds(p.Duration), nil
}

// fetchTransit 高德跨城公交：取首条方案的里程/时长，并判断是否含高铁/动车段。
func fetchTransit(client *http.Client, key, origin, dest, ocity, dcity string) (km, min *int, hasRail *bool, err error) {
	u := fmt.Sprintf("%s?origin=%s&destination=%s&city=%s&cityd=%s&key=%s",
		amapTransitURL, url.QueryEscape(origin), url.QueryEscape(dest),
		url.QueryEscape(ocity), url.QueryEscape(dcity), key)
	var body struct {
		Status string `json:"status"`
		Route  struct {
			Transits []struct {
				Duration string `json:"duration"`
				Distance string `json:"distance"`
				Segments []struct {
					// 非铁路段时高德返回空数组 []，铁路段为对象；用原文长度判非空
					Railway json.RawMessage `json:"railway"`
				} `json:"segments"`
			} `json:"transits"`
		} `json:"route"`
	}
	if err = getJSON(client, u, &body); err != nil {
		return nil, nil, nil, err
	}
	if body.Status != "1" || len(body.Route.Transits) == 0 {
		return nil, nil, nil, fmt.Errorf("amap transit: no route")
	}
	t := body.Route.Transits[0]
	rail := false
	for _, seg := range t.Segments {
		if len(seg.Railway) > 2 { // 排除 "[]" / "{}" 等空值
			rail = true
			break
		}
	}
	return kmFromMeters(t.Distance), minFromSeconds(t.Duration), &rail, nil
}

// getJSON 发起 GET 并解析 JSON；非 200 视为失败。
func getJSON(client *http.Client, u string, out interface{}) error {
	resp, err := client.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("amap http %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// kmFromMeters 米(字符串)→公里(四舍五入)。解析失败返回 nil。
func kmFromMeters(s string) *int {
	m, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	v := int(math.Round(m / 1000))
	return &v
}

// minFromSeconds 秒(字符串)→分钟(四舍五入)。解析失败返回 nil。
func minFromSeconds(s string) *int {
	sec, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	v := int(math.Round(sec / 60))
	return &v
}
