package scrapper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"

	"github.com/AFK068/bot/internal/infrastructure/logger"

	scrappertypes "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

type Service interface {
	PostTgChatID(ctx context.Context, id int64) error
	DeleteTgChatID(ctx context.Context, id int64) error
	PostLinks(ctx context.Context, tgChatID int64, link scrappertypes.AddLinkRequest) error
	DeleteLinks(ctx context.Context, tgChatID int64, link scrappertypes.RemoveLinkRequest) error
	GetLinks(ctx context.Context, tgChatID int64, tag ...string) (scrappertypes.ListLinksResponse, error)
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

func (c *Client) PostTgChatID(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/tg-chat/%d", c.BaseURL, id)
	c.Logger.Info("Posting TgChatID", "url", url, "id", id)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetHeader(echo.HeaderAccept, echo.MIMEApplicationJSON).
		Post(url)
	if err != nil {
		c.Logger.Error("Failed to post TgChatID", "error", err)
		return fmt.Errorf("failed to do request: %w", err)
	}

	return c.handleResponse(resp.StatusCode(), resp.Body())
}

func (c *Client) DeleteTgChatID(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/tg-chat/%d", c.BaseURL, id)
	c.Logger.Info("Deleting TgChatID", "url", url, "id", id)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetHeader(echo.HeaderAccept, echo.MIMEApplicationJSON).
		Delete(url)
	if err != nil {
		c.Logger.Error("Failed to delete TgChatID", "error", err)
		return fmt.Errorf("failed to do request: %w", err)
	}

	return c.handleResponse(resp.StatusCode(), resp.Body())
}

func (c *Client) PostLinks(ctx context.Context, tgChatID int64, link scrappertypes.AddLinkRequest) error {
	url := fmt.Sprintf("%s/links", c.BaseURL)
	c.Logger.Info("Posting Links", "url", url, "tgChatID", tgChatID, "link", link)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetHeader(echo.HeaderAccept, echo.MIMEApplicationJSON).
		SetHeader("Tg-Chat-Id", fmt.Sprintf("%d", tgChatID)).
		SetBody(link).
		Post(url)
	if err != nil {
		c.Logger.Error("Failed to post Links", "error", err)
		return fmt.Errorf("failed to do request: %w", err)
	}

	return c.handleResponse(resp.StatusCode(), resp.Body())
}

func (c *Client) DeleteLinks(ctx context.Context, tgChatID int64, link scrappertypes.RemoveLinkRequest) error {
	url := fmt.Sprintf("%s/links", c.BaseURL)
	c.Logger.Info("Deleting Links", "url", url, "tgChatID", tgChatID, "link", link)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetHeader(echo.HeaderAccept, echo.MIMEApplicationJSON).
		SetHeader("Tg-Chat-Id", fmt.Sprintf("%d", tgChatID)).
		SetBody(link).
		Delete(url)
	if err != nil {
		c.Logger.Error("Failed to delete Links", "error", err)
		return fmt.Errorf("failed to do request: %w", err)
	}

	return c.handleResponse(resp.StatusCode(), resp.Body())
}

func (c *Client) GetLinks(ctx context.Context, tgChatID int64, tag ...string) (scrappertypes.ListLinksResponse, error) {
	url := fmt.Sprintf("%s/links", c.BaseURL)
	c.Logger.Info("Getting Links", "url", url, "tgChatID", tgChatID)

	req := c.Client.R().
		SetContext(ctx).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetHeader(echo.HeaderAccept, echo.MIMEApplicationJSON).
		SetHeader("Tg-Chat-Id", fmt.Sprintf("%d", tgChatID))

	if len(tag) > 0 && tag[0] != "" {
		req.SetQueryParam("tag", tag[0])
	}

	resp, err := req.Get(url)
	if err != nil {
		c.Logger.Error("Failed to get Links", "error", err)
		return scrappertypes.ListLinksResponse{}, fmt.Errorf("failed to do request: %w", err)
	}

	if err := c.handleResponse(resp.StatusCode(), resp.Body()); err != nil {
		return scrappertypes.ListLinksResponse{}, err
	}

	var links scrappertypes.ListLinksResponse
	if err := json.Unmarshal(resp.Body(), &links); err != nil {
		return scrappertypes.ListLinksResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return links, nil
}
