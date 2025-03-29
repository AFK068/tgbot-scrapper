package github

import "time"

type repositoryDTO struct {
	ID          int64     `json:"id"`
	URL         string    `json:"url"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Owner       ownerDTO  `json:"owner"`
}

func (r *repositoryDTO) toRepository() *Repository {
	return NewRepository(r.ID, r.URL, r.UpdatedAt, r.CreatedAt, r.Description, r.Owner.Login)
}

// In GitHub terminology, a pull request is included in a request for issues.
type issueDTO struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAtAt time.Time `json:"created_at"`
	User        userDTO   `json:"user"`

	// The pull request and the issue are not explicitly separated in the requests,
	// so if any of these fields are not null it is of this type.
	PullRequest     *pullRequestDTO     `json:"pull_request"`
	SubIssueSummary *subIssueSummaryDTO `json:"sub_issue_summary"`
}

func (i *issueDTO) toIssue(issueType IssueType) *Issue {
	return NewIssue(issueType, i.ID, i.Title, i.Body, i.UpdatedAt, i.CreatedAtAt)
}

type ownerDTO struct {
	Login string `json:"login"`
}

type userDTO struct {
	Login string `json:"login"`
}

type pullRequestDTO struct {
	URL string `json:"url"`
}

type subIssueSummaryDTO struct {
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
}
