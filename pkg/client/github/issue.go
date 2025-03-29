package github

import "time"

type IssueType string

const (
	IssueTypePullRequest IssueType = "PullRequest"
	IssueTypeIssue       IssueType = "Issue"
)

type Issue struct {
	Type      IssueType
	ID        int64
	Title     string
	Body      string
	UpdatedAt time.Time
	CreatedAt time.Time
}

func NewIssue(issueType IssueType, id int64, title, body string, updatedAt, createdAt time.Time) *Issue {
	return &Issue{
		Type:      issueType,
		ID:        id,
		Title:     title,
		Body:      body,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}
}
