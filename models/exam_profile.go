package models

// ExamProfile 测档档案表（tbl_exam_profile）。一用户一行（user_id 唯一），整体 upsert。
// 保存用户选择的参考省份/首选科目/再选/分数或位次,供刷新或换设备登录后回填。
// Score/Rank 用指针,区分「填了 0」与「没填(NULL)」;二者互斥(填分则位次空,反之亦然)。
type ExamProfile struct {
	Meta
	UserID    int64     `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user;comment:用户ID(唯一)"`
	ScoreMode string    `json:"score_mode" gorm:"column:score_mode;type:varchar(16);not null;default:pre;comment:出分前pre/出分后post"`
	Province  string    `json:"province" gorm:"column:province;type:varchar(32);not null;default:'';comment:参考省份(空=未选)"`
	Subject   string    `json:"subject" gorm:"column:subject;type:varchar(16);not null;default:'';comment:首选科目(空=未选)"`
	Electives MajorList `json:"electives" gorm:"column:electives;type:json;comment:再选科目名数组"`
	Score     *int      `json:"score" gorm:"column:score;comment:高考分数(NULL=未填)"`
	Rank      *int      `json:"rank" gorm:"column:rank;comment:全省位次(NULL=未填)"`
}
