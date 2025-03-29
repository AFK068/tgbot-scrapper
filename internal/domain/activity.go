package domain

import (
	"time"

	bottypes "github.com/AFK068/bot/internal/api/openapi/bot/v1"
)

type ActivityType string

const (
	StackoverflowComment  ActivityType = "stackoverflow_comment"
	StackoverflowAnswer   ActivityType = "stackoverflow_answer"
	StackoverflowQuestion ActivityType = "stackoverflow_question"

	GitHubRepository  ActivityType = "github_repository"
	GitHubIssue       ActivityType = "github_issue"
	GitHubPullRequest ActivityType = "github_pull_request"
)

type Activity struct {
	Type      ActivityType
	Title     string
	CreatedAt time.Time
	Body      string
	UserName  string
}

func NewActivity(
	activityType ActivityType,
	title string,
	createdAt time.Time,
	body string,
	userName string,
) *Activity {
	return &Activity{
		Type:      activityType,
		Title:     title,
		CreatedAt: createdAt,
		Body:      body,
		UserName:  userName,
	}
}

func (a *Activity) MapActivityTypeToBotAPI() *bottypes.LinkUpdateType {
	switch a.Type {
	case StackoverflowComment:
		stackoverflowComment := bottypes.StackoverflowComment
		return &stackoverflowComment
	case StackoverflowAnswer:
		stackoverflowAnswer := bottypes.StackoverflowAnswer
		return &stackoverflowAnswer
	case StackoverflowQuestion:
		stackoverflowQuestion := bottypes.StackoverflowQuestion
		return &stackoverflowQuestion
	case GitHubRepository:
		githubRepository := bottypes.GithubRepository
		return &githubRepository
	case GitHubIssue:
		githubIssue := bottypes.GithubIssue
		return &githubIssue
	case GitHubPullRequest:
		githubPullRequest := bottypes.GithubPullRequest
		return &githubPullRequest
	}

	return nil
}
