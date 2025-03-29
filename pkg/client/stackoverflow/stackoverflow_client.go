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

// GetQuestion retrieves a question from the Stack Overflow API using its URL.
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

// GetActivity retrieves the activity of a given question since the last check time.
// It fetches the question's answers and comments from the Stack Overflow API
// and checks if they have been updated since the last check time. If so, it creates
// activity entries for each updated answer and comment and returns them.
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
	questionAnswerActivity, err := c.GetQuestionAnswerActivity(ctx, question, lastCheckTime)
	if err != nil {
		return nil, err
	}

	activities = append(activities, questionAnswerActivity...)

	// Get comments for the question.
	questionCommentActivity, err := c.GetQuestionCommentActivity(ctx, question, lastCheckTime)
	if err != nil {
		return nil, err
	}

	activities = append(activities, questionCommentActivity...)

	return activities, nil
}

// GetQuestionCommentActivity retrieves the activity of comments for a given question.
// It fetches the comments from the Stack Overflow API and checks if they have been
// updated since the last check time. If so, it creates activity entries for each
// updated comment and returns them.
func (c *Client) GetQuestionCommentActivity(ctx context.Context, question *Question, lastCheckTime time.Time) ([]*Activity, error) {
	var activities []*Activity

	commentURL := fmt.Sprintf("%s/questions/%d/comments?site=stackoverflow&filter=withbody", c.BaseURL, question.ID)

	commentItems, err := getItems[CommentResponse](ctx, c.Client, commentURL)
	if err != nil {
		return nil, err
	}

	if len(commentItems.Items) != 0 {
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
	}

	return activities, nil
}

// GetQuestionAnswerActivity retrieves the activity of answers for a given question.
// It fetches the answers from the Stack Overflow API and checks if they have been
// updated since the last check time. If so, it creates activity entries for each
// updated answer and returns them.
func (c *Client) GetQuestionAnswerActivity(ctx context.Context, question *Question, lastCheckTime time.Time) ([]*Activity, error) {
	var activities []*Activity

	answerURL := fmt.Sprintf("%s/questions/%d/answers?site=stackoverflow&filter=withbody", c.BaseURL, question.ID)

	answerItems, err := getItems[AnswerResponse](ctx, c.Client, answerURL)
	if err != nil {
		return nil, err
	}

	if len(answerItems.Items) != 0 {
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
	}

	return activities, nil
}

func getItems[T any](ctx context.Context, client *resty.Client, url string) (*T, error) {
	result := new(T)

	resp, err := client.R().
		SetContext(ctx).
		SetResult(result).
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
