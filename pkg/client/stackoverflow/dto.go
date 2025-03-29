package stackoverflow

type questionResponseDTO struct {
	Items []*questionDTO `json:"items"`
}

type answerResponseDTO struct {
	Items []*answerDTO `json:"items"`
}

type commentResponseDTO struct {
	Items []*commentDTO `json:"items"`
}

type questionDTO struct {
	ID               int64    `json:"question_id"`
	Owner            ownerDTO `json:"owner"`
	LastActivityDate int64    `json:"last_activity_date"`
	LastEditDate     int64    `json:"last_edit_date"`
	Tags             []string `json:"tags"`
	Body             string   `json:"body"`
}

func (q *questionDTO) toQuestion() *Question {
	return NewQuestion(
		q.ID,
		q.Owner.DisplayName,
		q.LastActivityDate,
		q.LastEditDate,
		q.Body,
		q.Tags)
}

type commentDTO struct {
	ID        int64    `json:"comment_id"`
	Owner     ownerDTO `json:"owner"`
	CreatedAt int64    `json:"creation_date"`
	Body      string   `json:"body"`
}

type answerDTO struct {
	ID               int64    `json:"answer_id"`
	Owner            ownerDTO `json:"owner"`
	Body             string   `json:"body"`
	LastActivityDate int64    `json:"last_activity_date"`
}

type ownerDTO struct {
	DisplayName string `json:"display_name"`
}
