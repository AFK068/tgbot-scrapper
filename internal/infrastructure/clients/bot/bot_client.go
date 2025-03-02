package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

type Service interface {
	PostUpdates(ctx context.Context, update api.LinkUpdate) error
}

type Client struct {
	BaseURL string
	Client  *resty.Client
	Logger  *logger.Logger
}

func NewClient(url string, log *logger.Logger) *Client {
	return &Client{
		Client:  resty.New(),
		BaseURL: url,
		Logger:  log,
	}
}

func (c *Client) PostUpdates(ctx context.Context, update api.LinkUpdate) error {
	url := fmt.Sprintf("%s/updates", c.BaseURL)

	c.Logger.Info("Sending update to URL: ", "url", url)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetHeader(echo.HeaderAccept, echo.MIMEApplicationJSON).
		SetBody(update).
		Post(url)
	if err != nil {
		c.Logger.Error("Failed to do request: ", "error", err)
		return fmt.Errorf("failed to do request: %w", err)
	}

	c.Logger.Info("Received response with status code: ", "status_code", resp.StatusCode())

	switch resp.StatusCode() {
	case http.StatusOK:
		c.Logger.Info("Update posted successfully")
		return nil
	case http.StatusBadRequest:
		var apiErr api.ApiErrorResponse
		if err := json.Unmarshal(resp.Body(), &apiErr); err != nil {
			c.Logger.Error("Failed to decode error response: ", "error", err)
			return fmt.Errorf("failed to decode error response: %w", err)
		}

		c.Logger.Error("Bad request: ", "description", *apiErr.Description)
		return fmt.Errorf("bad request: %s", *apiErr.Description)
	default:
		c.Logger.Error("Unexpected status code: ", "status_code", resp.StatusCode())
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
}
