package models

// Meta 所有表内嵌的公共字段。时间戳统一存 Unix 秒。
type Meta struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CreateTime int64 `json:"create_time" gorm:"column:create_time;autoCreateTime;comment:创建时间(Unix秒)"`
	UpdateTime int64 `json:"update_time" gorm:"column:update_time;autoUpdateTime;comment:更新时间(Unix秒)"`
}
