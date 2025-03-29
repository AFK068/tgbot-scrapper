package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	BaseGitHubAPIURL = "https://api.github.com"
	TrimBodyLimit    = 200
)

type Client struct {
	BaseURL string
	Client  *resty.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: BaseGitHubAPIURL,
		Client:  resty.New().SetTimeout(10 * time.Second),
	}
}

func (c *Client) GetRepo(ctx context.Context, questionURL string) (*Repository, error) {
	ownerName, repoName, err := getOwnerAndRepo(questionURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/repos/%s/%s", c.BaseURL, ownerName, repoName)

	var repoDTO repositoryDTO

	resp, err := c.Client.R().
		SetContext(ctx).
		SetResult(&repoDTO).
		Get(url)
	if err != nil {
		return nil, errors.New("failed to get repository")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("failed to get repository")
	}

	return repoDTO.toRepository(), nil
}

func (c *Client) GetActivity(ctx context.Context, repository *Repository, lastCheckTime time.Time) ([]*Activity, error) {
	var activities []*Activity

	if repository.UpdatedAt.After(lastCheckTime) {
		activities = append(activities, NewActivity(
			ActivityTypeRepository,
			repository.Description,
			repository.UpdatedAt,
			"",
			repository.Owner,
		))
	}

	page := 0

	for {
		page++

		issues, err := c.GetIssuesByPage(ctx, repository.URL, page)
		if err != nil {
			return nil, err
		}

		if len(issues) == 0 {
			break
		}

		for _, issue := range issues {
			if issue.UpdatedAt.After(lastCheckTime) {
				activities = append(activities, NewActivity(
					ActivityType(issue.Type),
					issue.Title,
					issue.UpdatedAt,
					issue.Body,
					repository.Owner,
				))
			}
		}
	}

	return activities, nil
}

func (c *Client) GetIssuesByPage(ctx context.Context, questionURL string, page int) ([]*Issue, error) {
	ownerName, repoName, err := getOwnerAndRepo(questionURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues?page=%d", c.BaseURL, ownerName, repoName, page)

	var issues []*issueDTO

	resp, err := c.Client.R().
		SetContext(ctx).
		SetResult(&issues).
		Get(url)
	if err != nil {
		return nil, errors.New("failed to get issues")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("failed to get issues")
	}

	var result []*Issue

	for _, issue := range issues {
		issue.Body = trimBody(issue.Body)

		if issue.PullRequest != nil {
			result = append(result, issue.toIssue(IssueTypePullRequest))
		} else {
			result = append(result, issue.toIssue(IssueTypeIssue))
		}
	}

	return result, nil
}

func trimBody(body string) string {
	if len(body) > TrimBodyLimit {
		return body[:TrimBodyLimit] + "..."
	}

	return body
}

func getOwnerAndRepo(url string) (owner, repo string, err error) {
	re := regexp.MustCompile(`(?i)(?:github\.com[/:])?([^/]+)/([^/]+?)(?:\.git)?$`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 3 {
		return "", "", errors.New("invalid GitHub repository URL format")
	}

	owner = strings.TrimSpace(matches[1])
	repo = strings.TrimSpace(strings.TrimSuffix(matches[2], ".git"))

	if owner == "" || repo == "" {
		return "", "", errors.New("empty owner or repository name")
	}

	return owner, repo, nil
}
