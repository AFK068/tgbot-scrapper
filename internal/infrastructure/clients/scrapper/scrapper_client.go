package scrapper

import (
	"context"
	"encoding/json"
	"fmt"

	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
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

func (c *Client) PostTgChatID(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/tg-chat/%d", c.BaseURL, id)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Post(url)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	return handleResponse(resp.StatusCode(), resp.Body())
}

func (c *Client) DeleteTgChatID(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/tg-chat/%d", c.BaseURL, id)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Delete(url)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	return handleResponse(resp.StatusCode(), resp.Body())
}

func (c *Client) PostLinks(ctx context.Context, tgChatID int64, link api.AddLinkRequest) error {
	url := fmt.Sprintf("%s/links", c.BaseURL)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Tg-Chat-Id", fmt.Sprintf("%d", tgChatID)).
		SetBody(link).
		Post(url)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	return handleResponse(resp.StatusCode(), resp.Body())
}

func (c *Client) DeleteLinks(ctx context.Context, tgChatID int64, link api.RemoveLinkRequest) error {
	url := fmt.Sprintf("%s/links", c.BaseURL)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Tg-Chat-Id", fmt.Sprintf("%d", tgChatID)).
		SetBody(link).
		Delete(url)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	return handleResponse(resp.StatusCode(), resp.Body())
}

func (c *Client) GetLinks(ctx context.Context, tgChatID int64) (api.ListLinksResponse, error) {
	url := fmt.Sprintf("%s/links", c.BaseURL)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Tg-Chat-Id", fmt.Sprintf("%d", tgChatID)).
		Get(url)
	if err != nil {
		return api.ListLinksResponse{}, fmt.Errorf("failed to do request: %w", err)
	}

	if err := handleResponse(resp.StatusCode(), resp.Body()); err != nil {
		return api.ListLinksResponse{}, err
	}

	var links api.ListLinksResponse
	if err := json.Unmarshal(resp.Body(), &links); err != nil {
		return api.ListLinksResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return links, nil
}
