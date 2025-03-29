package github_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AFK068/bot/pkg/client/github"
)

func Test_GetRepo_Success(t *testing.T) {
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

func Test_GetRepo_InvalidLink(t *testing.T) {
	client := github.NewClient()
	_, err := client.GetRepo(context.Background(), "https://bad_link")

	assert.Error(t, err)
}

func Test_GetActivity_Success(t *testing.T) {
	lastCheckTime := time.Now().Add(-244 * time.Hour)
	expectedTime := time.Now()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/test/test":
			response := map[string]interface{}{
				"updated_at":  expectedTime.Format(time.RFC3339),
				"description": "Test repo",
				"owner": map[string]interface{}{
					"login": "testuser",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(response)
			require.NoError(t, err)
		case "/repos/test/test/issues":
			page := r.URL.Query().Get("page")
			if page == "" || page == "1" {
				response := []map[string]interface{}{
					{
						"title":        "Test issue",
						"updated_at":   expectedTime.Format(time.RFC3339),
						"body":         "Test issue body",
						"pull_request": nil,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(w).Encode(response)
				require.NoError(t, err)
			} else {
				w.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(w).Encode([]interface{}{})
				require.NoError(t, err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	defer server.Close()

	client := github.NewClient()
	client.BaseURL = server.URL
	client.Client = client.Client.SetBaseURL(server.URL)

	repo := &github.Repository{
		URL:         "https://github.com/test/test",
		Description: "Test repo",
		UpdatedAt:   expectedTime,
		Owner:       "testuser",
	}

	activities, err := client.GetActivity(context.Background(), repo, lastCheckTime)
	require.NoError(t, err)
	assert.Len(t, activities, 2)
}

func Test_GetIssuesByPage_Success(t *testing.T) {
	expectedTime := time.Now()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/test/test/issues", r.URL.Path)
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		response := []map[string]interface{}{
			{
				"title":        "Test issue",
				"updated_at":   expectedTime.Format(time.RFC3339),
				"body":         "Test issue body",
				"pull_request": nil,
			},
			{
				"title":        "Test PR",
				"updated_at":   expectedTime.Format(time.RFC3339),
				"body":         "Test PR body",
				"pull_request": map[string]interface{}{},
			},
		}

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))

	defer server.Close()

	client := github.NewClient()
	client.BaseURL = server.URL
	client.Client = client.Client.SetBaseURL(server.URL)

	issues, err := client.GetIssuesByPage(context.Background(), "https://github.com/test/test", 1)
	require.NoError(t, err)
	require.Len(t, issues, 2)
	assert.Equal(t, "Test issue", issues[0].Title)
	assert.Equal(t, github.IssueTypeIssue, issues[0].Type)
	assert.Equal(t, "Test PR", issues[1].Title)
	assert.Equal(t, github.IssueTypePullRequest, issues[1].Type)
}
