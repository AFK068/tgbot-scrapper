package github

import "time"

type ActivityType string

const (
	ActivityTypePullRequest ActivityType = "PullRequest"
	ActivityTypeIssue       ActivityType = "Issue"
	ActivityTypeRepository  ActivityType = "Repository"
)

type Activity struct {
	Type      ActivityType
	Title     string
	CreatedAt time.Time
	Body      string
	UserName  string
}

func NewActivity(activityType ActivityType, title string, createdAt time.Time, body, userName string) *Activity {
	return &Activity{
		Type:      activityType,
		Title:     title,
		CreatedAt: createdAt,
		Body:      body,
		UserName:  userName,
	}
}
