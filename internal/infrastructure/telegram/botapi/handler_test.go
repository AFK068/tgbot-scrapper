package botapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/internal/infrastructure/telegram/botapi"

	bottypes "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	botmocks "github.com/AFK068/bot/internal/application/bot/mocks"
)

func Test_PostUpdates_Success(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := botapi.NewBotHandler(botMock, logger.NewDiscardLogger())

	botMock.On("SendMessage", int64(123), "Link updated: https://test\nDescription: Test description").Once()
	botMock.On("SendMessage", int64(456), "Link updated: https://test").Once()

	testCases := []struct {
		name string
		body bottypes.LinkUpdate
	}{
		{
			name: "With description",
			body: bottypes.LinkUpdate{
				TgChatIds:   &[]int64{123},
				Url:         aws.String("https://test"),
				Description: aws.String("Test description"),
			},
		},
		{
			name: "Without description",
			body: bottypes.LinkUpdate{
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

func Test_PostUpdates_InvalidBody(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := botapi.NewBotHandler(botMock, logger.NewDiscardLogger())

	req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader([]byte(`Invalid_body`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.PostUpdates(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	botMock.AssertNotCalled(t, "SendMessage")
}

func Test_PostUpdates_EmptyTgChatIDs(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := botapi.NewBotHandler(botMock, logger.NewDiscardLogger())

	testCases := []struct {
		name string
		body bottypes.LinkUpdate
	}{
		{
			name: "Empty TgChatIDs",
			body: bottypes.LinkUpdate{
				TgChatIds: &[]int64{},
				Url:       aws.String("https://test"),
			},
		},
		{
			name: "Nil TgChatIDs",
			body: bottypes.LinkUpdate{
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

func Test_PostUpdates_EmptyURL(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := botapi.NewBotHandler(botMock, logger.NewDiscardLogger())

	testCases := []struct {
		name string
		body bottypes.LinkUpdate
	}{
		{
			name: "Empty URL",
			body: bottypes.LinkUpdate{
				TgChatIds: &[]int64{123},
				Url:       aws.String(""),
			},
		},
		{
			name: "Nil URL",
			body: bottypes.LinkUpdate{
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

func Test_PostUpdates_EmptyDescription(t *testing.T) {
	botMock := botmocks.NewService(t)
	h := botapi.NewBotHandler(botMock, logger.NewDiscardLogger())

	botMock.On("SendMessage", int64(123), "Link updated: https://test").Once()

	reqBody := bottypes.LinkUpdate{
		TgChatIds: &[]int64{123},
		Url:       aws.String("https://test"),
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.PostUpdates(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	botMock.AssertExpectations(t)
}
