package stackoverflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

const (
	BaseStackOverflowAPIURL = "https://api.stackexchange.com/2.2"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: BaseStackOverflowAPIURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetQuestion(ctx context.Context, questionURL string) (*Question, error) {
	questionID, err := getIDFromURL(questionURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/questions/%s?site=stackoverflow", c.BaseURL, questionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get question")
	}

	var response QuestionResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, errors.New("question not found")
	}

	return &response.Items[0], nil
}

func getIDFromURL(url string) (string, error) {
	reg := regexp.MustCompile(`questions/(\d+)`)

	matches := reg.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", errors.New("invalid question url")
	}

	return matches[1], nil
}
