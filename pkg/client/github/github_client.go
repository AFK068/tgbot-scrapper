package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	BaseGitHubAPIURL = "https://api.github.com"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: BaseGitHubAPIURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetRepo(ctx context.Context, url string) (*Repository, error) {
	ownerName, repoName, err := getOwnerAndRepo(url)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/repos/%s/%s", c.BaseURL, ownerName, repoName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get repository")
	}

	var repo Repository
	if err := json.NewDecoder(res.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
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
