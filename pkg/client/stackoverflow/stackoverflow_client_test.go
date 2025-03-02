package stackoverflow_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AFK068/bot/pkg/client/stackoverflow"
)

func TestGetRepo_Success(t *testing.T) {
	expectedTime := time.Unix(123456789, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/questions/123", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"question_id":        123,
					"last_activity_date": 123456789,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))

	defer server.Close()

	client := stackoverflow.NewClient()
	client.BaseURL = server.URL
	client.Client = client.Client.SetBaseURL(server.URL)

	question, err := client.GetQuestion(context.Background(), "https://stackoverflow.com/questions/123")
	assert.NoError(t, err)

	assert.Equal(t, int64(123), question.QuestionID)
	assert.Equal(t, expectedTime.Unix(), question.LastActivityDate)
}

func TestGetRepo_InvalidLink(t *testing.T) {
	client := stackoverflow.NewClient()
	_, err := client.GetQuestion(context.Background(), "https://bad_link")

	assert.Error(t, err)
}
