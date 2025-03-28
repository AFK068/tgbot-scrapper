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

func Test_GetRepo_Success(t *testing.T) {
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

	assert.Equal(t, int64(123), question.ID)
	assert.Equal(t, expectedTime.Unix(), question.LastActivityDate)
}

func Test_GetRepo_InvalidLink(t *testing.T) {
	client := stackoverflow.NewClient()
	_, err := client.GetQuestion(context.Background(), "https://bad_link")

	assert.Error(t, err)
}

func Test_GetActivity_Success(t *testing.T) {
	answerResponse := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"answer_id":          1,
				"body":               "Test answer body",
				"last_activity_date": 100,
				"owner": map[string]interface{}{
					"display_name": "AnswerUser",
				},
			},
		},
	}

	commentResponse := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"comment_id":    101,
				"creation_date": 200,
				"body":          "Test comment body",
				"owner": map[string]interface{}{
					"display_name": "CommentUser",
				},
			},
		},
	}

	answerBytes, err := json.Marshal(answerResponse)
	assert.NoError(t, err)

	commentBytes, err := json.Marshal(commentResponse)
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/questions/123/answers":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write(answerBytes)
			assert.NoError(t, err)
		case "/questions/123/comments":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write(commentBytes)
			assert.NoError(t, err)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	defer server.Close()

	client := stackoverflow.NewClient()
	client.BaseURL = server.URL

	question := &stackoverflow.Question{
		ID:           123,
		LastEditDate: 40,
		Body:         "Test question body",
		Tags:         []string{"go", "api"},
		Owner:        stackoverflow.Owner{DisplayName: "QuestionUser"},
	}

	lastCheckTime := time.Unix(50, 0)

	activities, err := client.GetActivity(context.Background(), question, lastCheckTime)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(activities), 2)
}
