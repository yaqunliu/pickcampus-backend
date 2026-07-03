package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// MajorList 候选专业名数组，落库为 JSON 文本。空数组存 "[]"，不存 NULL。
type MajorList []string

// Value 实现 driver.Valuer：写库时序列化为 JSON。
func (m MajorList) Value() (driver.Value, error) {
	if m == nil {
		return "[]", nil
	}
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner：读库时从 JSON 反序列化。
func (m *MajorList) Scan(src interface{}) error {
	if src == nil {
		*m = MajorList{}
		return nil
	}
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal([]byte(v), m)
	default:
		return errors.New("MajorList.Scan: 不支持的类型")
	}
}

// Candidate 候选院校表（tbl_candidate）。一行 = 某用户候选某院校。
// school_id 是前端静态院校库的字符串 id，后端无院校表，故存 varchar、不做外键。
// majors 存该校被标注的候选专业名数组（JSON），删行即级联清除。
type Candidate struct {
	Meta
	UserID   int64     `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user_school,priority:1;index:idx_user;comment:用户ID"`
	SchoolID string    `json:"school_id" gorm:"column:school_id;type:varchar(64);not null;uniqueIndex:idx_user_school,priority:2;comment:院校ID(前端静态库字符串id)"`
	Majors   MajorList `json:"majors" gorm:"column:majors;type:json;comment:候选专业名数组"`
}
