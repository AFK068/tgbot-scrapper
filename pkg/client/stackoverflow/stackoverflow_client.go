package stackoverflow

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	BaseStackOverflowAPIURL = "https://api.stackexchange.com/2.2"
	TrimBodyLimit           = 200
)

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

	url := fmt.Sprintf("%s/questions/%s?site=stackoverflow&filter=withbody", c.BaseURL, questionID)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetResult(&QuestionResponse{}).
		Get(url)
	if err != nil {
		return nil, ErrFailedToGetQuestion
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, ErrFailedToGetQuestion
	}

	quesion := resp.Result().(*QuestionResponse)
	if len(quesion.Items) == 0 {
		return nil, ErrQuestionNotFound
	}

	// Trim the body of the question to a certain limit.
	quesion.Items[0].Body = trimBody(quesion.Items[0].Body)

	return quesion.Items[0], nil
}

func (c *Client) GetActivity(ctx context.Context, question *Question, lastCheckTime time.Time) ([]*Activity, error) {
	var activities []*Activity

	// Check main question activity.
	if question.LastEditDate > lastCheckTime.Unix() {
		activity := NewActivity(
			ActivityTypeQuestion,
			question.LastEditDate,
			trimBody(question.Body),
			question.Tags,
			question.Owner.DisplayName,
		)

		activities = append(activities, activity)
	}

	// Get answers for the question.
	answerURL := fmt.Sprintf("%s/questions/%d/answers?site=stackoverflow&filter=withbody", c.BaseURL, question.ID)

	answerItems, err := getItems[AnswerResponse](ctx, answerURL)
	if err != nil {
		return nil, err
	}

	if len(answerItems.Items) == 0 {
		return nil, ErrNoAnswersFound
	}

	for _, answer := range answerItems.Items {
		if answer.LastActivityDate > lastCheckTime.Unix() {
			activity := NewActivity(
				ActivityTypeAnswer,
				answer.LastActivityDate,
				trimBody(answer.Body),
				question.Tags,
				answer.Owner.DisplayName,
			)

			activities = append(activities, activity)
		}
	}

	// Get comments for the question.
	commentURL := fmt.Sprintf("%s/questions/%d/comments?site=stackoverflow&filter=withbody", c.BaseURL, question.ID)

	commentItems, err := getItems[CommentResponse](ctx, commentURL)
	if err != nil {
		return nil, err
	}

	if len(answerItems.Items) == 0 {
		return nil, ErrNoAnswersFound
	}

	for _, comment := range commentItems.Items {
		if comment.CreatedAt > lastCheckTime.Unix() {
			activity := NewActivity(
				ActivityTypeComment,
				comment.CreatedAt,
				trimBody(comment.Body),
				question.Tags,
				comment.Owner.DisplayName,
			)

			activities = append(activities, activity)
		}
	}

	return activities, nil
}

func getItems[T any](ctx context.Context, url string) (*T, error) {
	result := new(T)

	resp, err := resty.New().R().
		SetContext(ctx).
		SetResult(&result).
		Get(url)
	if err != nil {
		return result, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, ErrFailedToGetItems
	}

	return result, nil
}

func trimBody(body string) string {
	if len(body) > TrimBodyLimit {
		return body[:TrimBodyLimit] + "..."
	}

	return body
}

func getIDFromURL(url string) (string, error) {
	reg := regexp.MustCompile(`questions/(\d+)`)

	matches := reg.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", ErrInvalidQuestionURL
	}

	return matches[1], nil
}
