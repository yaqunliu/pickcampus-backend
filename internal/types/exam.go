package types

// ExamProfileDTO 测档档案对外结构（GET 返回 / PUT 入参共用）。
// score/rank 用指针,区分「未填(null)」与「填了 0」。province/subject 空串=未选。
type ExamProfileDTO struct {
	ScoreMode string   `json:"score_mode"`
	Province  string   `json:"province"`
	Subject   string   `json:"subject"`
	Electives []string `json:"electives"`
	Score     *int     `json:"score"`
	Rank      *int     `json:"rank"`
}
