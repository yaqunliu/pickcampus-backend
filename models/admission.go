package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// SubRating 大类招生下某具体专业的软科评级。
type SubRating struct {
	Major  string `json:"major"`
	Rating string `json:"rating"`
}

// SubRatingList 大类子专业评级数组，落库为 JSON 文本。空则存 "[]"。
type SubRatingList []SubRating

// Value 写库时序列化为 JSON。
func (s SubRatingList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan 读库时从 JSON 反序列化。
func (s *SubRatingList) Scan(src interface{}) error {
	if src == nil {
		*s = SubRatingList{}
		return nil
	}
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("SubRatingList.Scan: 不支持的类型")
	}
}

// Admission 录取记录表（tbl_admission）。院校级 + 专业级共用一张表：
//
//	Major == nil  → 院校级记录（某校在某省某科的整体录取线）
//	Major != nil  → 专业级记录（具体专业的录取线）
//
// 对应前端 AdmissionRecord（major 可选，major===undefined 区分两级）。
// university_id 是前端静态院校库的字符串 id，后端无院校表，存 varchar、不做外键。
type Admission struct {
	Meta
	// 索引精简为 2 个,匹配阶段 2 的实际查询:
	//   idx_province(province)      —— 按省取院校级/专业级列表(major 空/非空在此基础上过滤)
	//   idx_uni_major(university_id, major) —— 按校取去重开设专业(覆盖索引)
	UniversityID string  `json:"university_id" gorm:"column:university_id;type:varchar(64);not null;index:idx_uni_major,priority:1;comment:院校ID(前端静态库字符串id)"`
	Province     string  `json:"province" gorm:"column:province;type:varchar(32);not null;index:idx_province;comment:招生省份"`
	Subject      string  `json:"subject" gorm:"column:subject;type:varchar(16);not null;comment:科类(物理类/历史类/综合/文科/理科)"`
	Year         int     `json:"year" gorm:"column:year;not null;comment:年份"`
	MinRank      int     `json:"min_rank" gorm:"column:min_rank;not null;comment:最低录取位次"`
	MinScore     *int    `json:"min_score" gorm:"column:min_score;comment:最低录取分(NULL=缺)"`
	Source       string  `json:"source" gorm:"column:source;type:varchar(128);not null;default:'';comment:数据来源标注"`
	Major        *string `json:"major" gorm:"column:major;type:varchar(128);index:idx_uni_major,priority:2;comment:专业名(NULL=院校级记录)"`
	MajorCode    string  `json:"major_code" gorm:"column:major_code;type:varchar(32);not null;default:'';comment:专业代码"`
	ElectiveReq  string  `json:"elective_req" gorm:"column:elective_req;type:varchar(64);not null;default:'';comment:选科要求"`
	Batch        string  `json:"batch" gorm:"column:batch;type:varchar(32);not null;default:'';comment:录取批次"`
	Tuition      *int    `json:"tuition" gorm:"column:tuition;comment:学费(元/年,NULL=缺)"`
	Duration     string  `json:"duration" gorm:"column:duration;type:varchar(16);not null;default:'';comment:学制"`
	Rating       string  `json:"rating" gorm:"column:rating;type:varchar(8);not null;default:'';comment:软科专业评级"`
	RatingRank   *int    `json:"rating_rank" gorm:"column:rating_rank;comment:软科专业全国排名(NULL=缺)"`
	SubRatings   SubRatingList `json:"sub_ratings" gorm:"column:sub_ratings;type:json;comment:大类子专业评级数组"`
}
