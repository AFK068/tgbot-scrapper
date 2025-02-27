package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	botmocks "github.com/AFK068/bot/internal/application/bot/mocks"
	handler "github.com/AFK068/bot/internal/infrastructure/handler/bot"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostUpdates_Success(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := handler.NewBotHandler(botMock)

	botMock.On("SendMessage", int64(123), "Link updated: https://test\nDescription: Test description").Once()
	botMock.On("SendMessage", int64(456), "Link updated: https://test").Once()

	testCases := []struct {
		name string
		body botapi.LinkUpdate
	}{
		{
			name: "With description",
			body: botapi.LinkUpdate{
				TgChatIds:   &[]int64{123},
				Url:         aws.String("https://test"),
				Description: aws.String("Test description"),
			},
		},
		{
			name: "Without description",
			body: botapi.LinkUpdate{
				TgChatIds: &[]int64{456},
				Url:       aws.String("https://test"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			err = h.PostUpdates(c)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestPostUpdates_InvalidBody(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := handler.NewBotHandler(botMock)

	req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader([]byte(`Invalid_body`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.PostUpdates(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	botMock.AssertNotCalled(t, "SendMessage")
}

func TestPostUpdates_EmptyTgChatIDs(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := handler.NewBotHandler(botMock)

	testCases := []struct {
		name string
		body botapi.LinkUpdate
	}{
		{
			name: "Empty TgChatIDs",
			body: botapi.LinkUpdate{
				TgChatIds: &[]int64{},
				Url:       aws.String("https://test"),
			},
		},
		{
			name: "Nil TgChatIDs",
			body: botapi.LinkUpdate{
				Url: aws.String("https://test"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			err = h.PostUpdates(c)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestPostUpdates_EmptyURL(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := handler.NewBotHandler(botMock)

	testCases := []struct {
		name string
		body botapi.LinkUpdate
	}{
		{
			name: "Empty URL",
			body: botapi.LinkUpdate{
				TgChatIds: &[]int64{123},
				Url:       aws.String(""),
			},
		},
		{
			name: "Nil URL",
			body: botapi.LinkUpdate{
				TgChatIds: &[]int64{123},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			err = h.PostUpdates(c)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}
