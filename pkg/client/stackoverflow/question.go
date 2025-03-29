package stackoverflow

type Question struct {
	ID               int64
	Name             string
	LastActivityDate int64
	LastEditDate     int64
	Tags             []string
	Body             string
}

func NewQuestion(id int64, name string, lastActivityDate, lastEditDate int64, body string, tags []string) *Question {
	return &Question{
		ID:               id,
		Name:             name,
		LastActivityDate: lastActivityDate,
		LastEditDate:     lastEditDate,
		Tags:             tags,
		Body:             body,
	}
}
