package scrapper_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	"github.com/AFK068/bot/internal/infrastructure/clients/scrapper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestPostTgChatID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		expectedPath := "/tg-chat/123"
		assert.Equal(t, expectedPath, r.URL.Path)

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL)
	err := client.PostTgChatID(context.Background(), 123)
	assert.NoError(t, err)
}

func TestDeleteTgChatID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		expectedPath := "/tg-chat/123"
		assert.Equal(t, expectedPath, r.URL.Path)

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL)
	err := client.DeleteTgChatID(context.Background(), 123)
	assert.NoError(t, err)
}

func TestPostLinks(t *testing.T) {
	reqBody := api.AddLinkRequest{
		Link:    aws.String("https://example.com"),
		Tags:    &[]string{"tag"},
		Filters: &[]string{"filter"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		expectedPath := "/links"
		assert.Equal(t, expectedPath, r.URL.Path)

		tgHeader := "Tg-Chat-ID"
		assert.Equal(t, r.Header.Get(tgHeader), "123")

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		var body api.AddLinkRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, reqBody, body)

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL)
	err := client.PostLinks(context.Background(), 123, reqBody)
	assert.NoError(t, err)
}

func TestDeleteLinks(t *testing.T) {
	reqBody := api.RemoveLinkRequest{
		Link: aws.String("https://example.com"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		expectedPath := "/links"
		assert.Equal(t, expectedPath, r.URL.Path)

		tgHeader := "Tg-Chat-ID"
		assert.Equal(t, r.Header.Get(tgHeader), "123")

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		var body api.RemoveLinkRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, reqBody, body)

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL)
	err := client.DeleteLinks(context.Background(), 123, reqBody)
	assert.NoError(t, err)
}

func TestGetLinks(t *testing.T) {
	response := api.ListLinksResponse{
		Links: &[]api.LinkResponse{
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

		expectedPath := "/links"
		assert.Equal(t, expectedPath, r.URL.Path)

		tgHeader := "Tg-Chat-ID"
		assert.Equal(t, r.Header.Get(tgHeader), "123")

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		resp, err := json.Marshal(response)
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		assert.NoError(t, err)
	}))

	defer server.Close()

	client := scrapper.NewClient(server.URL)
	resp, err := client.GetLinks(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, response, resp)
}
