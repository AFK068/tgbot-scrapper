package github

import "time"

type Repository struct {
	ID        int64     `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}
