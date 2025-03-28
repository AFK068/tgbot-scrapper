package stackoverflow

type QuestionResponse struct {
	Items []*Question `json:"items"`
}

type AnswerResponse struct {
	Items []*Answer `json:"items"`
}

type CommentResponse struct {
	Items []*Comment `json:"items"`
}

type Question struct {
	ID               int64    `json:"question_id"`
	Owner            Owner    `json:"owner"`
	LastActivityDate int64    `json:"last_activity_date"`
	LastEditDate     int64    `json:"last_edit_date"`
	Tags             []string `json:"tags"`
	Body             string   `json:"body"`
}

type Comment struct {
	ID        int64  `json:"comment_id"`
	Owner     Owner  `json:"owner"`
	CreatedAt int64  `json:"creation_date"`
	Body      string `json:"body"`
}

type Answer struct {
	ID               int64  `json:"answer_id"`
	Owner            Owner  `json:"owner"`
	Body             string `json:"body"`
	LastActivityDate int64  `json:"last_activity_date"`
}

type Owner struct {
	DisplayName string `json:"display_name"`
}
