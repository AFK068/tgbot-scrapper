package github_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AFK068/bot/pkg/client/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRepo_Success(t *testing.T) {
	expectedTime, err := time.Parse(time.RFC3339, "2011-01-26T19:14:43Z")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/test/test", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := map[string]interface{}{
			"id":         123,
			"updated_at": "2011-01-26T19:14:43Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))

	defer server.Close()

	client := github.NewClient()
	client.BaseURL = server.URL
	client.Client = client.Client.SetBaseURL(server.URL)

	repo, err := client.GetRepo(context.Background(), "https://github.com/test/test")

	require.NoError(t, err)
	assert.Equal(t, int64(123), repo.ID)
	assert.Equal(t, expectedTime, repo.UpdatedAt)
}

func TestGetRepo_InvalidLink(t *testing.T) {
	client := github.NewClient()
	_, err := client.GetRepo(context.Background(), "https://bad_link")

	assert.Error(t, err)
}
