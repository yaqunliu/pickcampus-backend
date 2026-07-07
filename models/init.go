package models

// AllTables 需要自动迁移的表清单。新增表在此追加。
var AllTables = []interface{}{
	&User{},
	&Candidate{},
	&ExamProfile{},
}
