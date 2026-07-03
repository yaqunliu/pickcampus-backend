package types

// AddCandidateRequest 加入候选院校。
type AddCandidateRequest struct {
	SchoolID string `json:"school_id" binding:"required,max=64"`
}

// CandidateMajorRequest 加/删候选专业（major 在 body，school_id 在 path）。
type CandidateMajorRequest struct {
	Major string `json:"major" binding:"required,max=128"`
}

// CandidateItem 单条候选（对外）。
type CandidateItem struct {
	SchoolID   string   `json:"school_id"`
	Majors     []string `json:"majors"`
	CreateTime int64    `json:"create_time"`
}
