package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	BaseURL string
	Client  *resty.Client
}

func NewClient(url string) *Client {
	return &Client{
		Client:  resty.New(),
		BaseURL: url,
	}
}

func (c *Client) PostUpdates(ctx context.Context, update api.LinkUpdate) error {
	url := fmt.Sprintf("%s/updates", c.BaseURL)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(update).
		Post(url)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		var apiErr api.ApiErrorResponse
		if err := json.Unmarshal(resp.Body(), &apiErr); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}

		return fmt.Errorf("bad request: %s", *apiErr.Description)
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
}
