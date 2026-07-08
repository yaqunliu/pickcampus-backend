package types

// SubRatingDTO 大类子专业评级(对外)。
type SubRatingDTO struct {
	Major  string `json:"major"`
	Rating string `json:"rating"`
}

// AdmissionRecordDTO 录取记录(对外)。
// 字段用 camelCase,与前端 AdmissionRecord / 原 public JSON 完全一致——
// 使前端 loader 只需换请求地址,AdmissionProvider 与组件零改动。
// omitempty:院校级记录只带基础字段(与原院校级 JSON 一致),专业级字段缺省时不下发。
type AdmissionRecordDTO struct {
	UniversityID string         `json:"universityId"`
	Province     string         `json:"province"`
	Subject      string         `json:"subject"`
	Year         int            `json:"year"`
	MinRank      int            `json:"minRank"`
	MinScore     *int           `json:"minScore,omitempty"`
	Source       string         `json:"source"`
	Major        *string        `json:"major,omitempty"`
	MajorCode    string         `json:"majorCode,omitempty"`
	ElectiveReq  string         `json:"electiveReq,omitempty"`
	Batch        string         `json:"batch,omitempty"`
	Tuition      *int           `json:"tuition,omitempty"`
	Duration     string         `json:"duration,omitempty"`
	Rating       string         `json:"rating,omitempty"`
	RatingRank   *int           `json:"ratingRank,omitempty"`
	SubRatings   []SubRatingDTO `json:"subRatings,omitempty"`
}
