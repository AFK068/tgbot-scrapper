package stackoverflow

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	BaseStackOverflowAPIURL = "https://api.stackexchange.com/2.2"
)

type QuestionFetcher interface {
	GetQuestion(ctx context.Context, questionURL string) (*Question, error)
}

type Client struct {
	BaseURL string
	Client  *resty.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: BaseStackOverflowAPIURL,
		Client:  resty.New().SetTimeout(10 * time.Second),
	}
}

func (c *Client) GetQuestion(ctx context.Context, questionURL string) (*Question, error) {
	questionID, err := getIDFromURL(questionURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/questions/%s?site=stackoverflow", c.BaseURL, questionID)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetResult(&QuestionResponse{}).
		Get(url)
	if err != nil {
		return nil, errors.New("failed to get question")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("failed to get question")
	}

	response := resp.Result().(*QuestionResponse)
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
