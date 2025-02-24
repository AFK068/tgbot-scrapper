package stackoverflow

type Question struct {
	QuestionID       int64 `json:"question_id"`
	LastActivityDate int64 `json:"last_activity_date"`
}

type QuestionResponse struct {
	Items []Question `json:"items"`
}
