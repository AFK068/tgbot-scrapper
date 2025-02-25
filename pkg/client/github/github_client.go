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

	resp, err := c.Client.R().
		SetContext(ctx).
		SetResult(&Repository{}).
		Get(url)
	if err != nil {
		return nil, errors.New("failed to get repository")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("failed to get repository")
	}

	return resp.Result().(*Repository), nil
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
