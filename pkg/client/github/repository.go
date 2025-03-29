package github

import "time"

type Repository struct {
	ID          int64
	URL         string
	UpdatedAt   time.Time
	CreatedAt   time.Time
	Description string
	Owner       string
}

func NewRepository(id int64, url string, updatedAt, createdAt time.Time, description, owner string) *Repository {
	return &Repository{
		ID:          id,
		URL:         url,
		UpdatedAt:   updatedAt,
		CreatedAt:   createdAt,
		Description: description,
		Owner:       owner,
	}
}
