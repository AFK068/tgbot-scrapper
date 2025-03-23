package scrapper_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/AFK068/bot/internal/infrastructure/clients/scrapper"
	"github.com/AFK068/bot/internal/infrastructure/logger"

	scrappertypes "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

func Test_PostTgChatID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		assert.Equal(t, "/tg-chat/123", r.URL.Path)

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL, logger.NewDiscardLogger())
	err := client.PostTgChatID(context.Background(), 123)
	assert.NoError(t, err)
}

func Test_DeleteTgChatID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		assert.Equal(t, "/tg-chat/123", r.URL.Path)

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL, logger.NewDiscardLogger())
	err := client.DeleteTgChatID(context.Background(), 123)
	assert.NoError(t, err)
}

func Test_PostLinks(t *testing.T) {
	reqBody := scrappertypes.AddLinkRequest{
		Link:    aws.String("https://example.com"),
		Tags:    &[]string{"tag"},
		Filters: &[]string{"filter"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		assert.Equal(t, "/links", r.URL.Path)

		assert.Equal(t, r.Header.Get("Tg-Chat-ID"), "123")

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		var body scrappertypes.AddLinkRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, reqBody, body)

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL, logger.NewDiscardLogger())
	err := client.PostLinks(context.Background(), 123, reqBody)
	assert.NoError(t, err)
}

func Test_DeleteLinks(t *testing.T) {
	reqBody := scrappertypes.RemoveLinkRequest{
		Link: aws.String("https://example.com"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		assert.Equal(t, "/links", r.URL.Path)

		assert.Equal(t, r.Header.Get("Tg-Chat-ID"), "123")

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		var body scrappertypes.RemoveLinkRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, reqBody, body)

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL, logger.NewDiscardLogger())
	err := client.DeleteLinks(context.Background(), 123, reqBody)
	assert.NoError(t, err)
}

func Test_GetLinks(t *testing.T) {
	response := scrappertypes.ListLinksResponse{
		Links: &[]scrappertypes.LinkResponse{
			{
				Url:     aws.String("https://example.com"),
				Tags:    &[]string{"tag"},
				Filters: &[]string{"filter"},
			},
		},
		Size: aws.Int32(1),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		assert.Equal(t, "/links", r.URL.Path)

		assert.Equal(t, r.Header.Get("Tg-Chat-ID"), "123")

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		resp, err := json.Marshal(response)
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		assert.NoError(t, err)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL, logger.NewDiscardLogger())
	resp, err := client.GetLinks(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, response, resp)
}
