package bot_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/AFK068/bot/internal/infrastructure/clients/bot"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestPostUpdates(t *testing.T) {
	reqBody := api.LinkUpdate{
		Url:         aws.String("https://example.com"),
		TgChatIds:   &[]int64{1, 2, 3},
		Description: aws.String("description"),
		Id:          aws.Int64(1),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		assert.Equal(t, "/updates", r.URL.Path)

		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept"), "application/json")

		var body api.LinkUpdate
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, reqBody, body)

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := bot.NewClient(server.URL)
	err := client.PostUpdates(context.Background(), reqBody)
	assert.NoError(t, err)
}
