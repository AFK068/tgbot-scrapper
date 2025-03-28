package domain

import "time"

var (
	StackoverflowType = "stackoverflow"
	GithubType        = "github"
)

type Link struct {
	UserAddID int64
	URL       string
	Type      string
	Tags      []string
	Filters   []string
	LastCheck time.Time
}
