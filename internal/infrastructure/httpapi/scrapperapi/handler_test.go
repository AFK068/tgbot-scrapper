package scrapperapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
	"github.com/AFK068/bot/internal/infrastructure/httpapi/scrapperapi"
	"github.com/AFK068/bot/internal/infrastructure/logger"

	scrappertypes "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	repomock "github.com/AFK068/bot/internal/domain/mocks"
)

func Test_PostTgChatId_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	repoMock.On("CheckUserExistence", mock.Anything, int64(123)).Return(false, nil)
	repoMock.On("RegisterChat", mock.Anything, int64(123)).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.PostTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_PostTgChatId_AlreadyExists(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	repoMock.On("CheckUserExistence", mock.Anything, int64(123)).Return(true, nil)

	req := httptest.NewRequest(http.MethodPost, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.PostTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_PostTgChatId_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	repoMock.On("CheckUserExistence", mock.Anything, int64(123)).Return(false, nil)
	repoMock.On("RegisterChat", mock.Anything, int64(123)).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.PostTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_DeleteTgChatId_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	repoMock.On("CheckUserExistence", mock.Anything, int64(123)).Return(true, nil)
	repoMock.On("DeleteChat", mock.Anything, int64(123)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.DeleteTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_DeleteTgChatId_UserNotFound(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	repoMock.On("CheckUserExistence", mock.Anything, int64(123)).Return(false, nil)

	req := httptest.NewRequest(http.MethodDelete, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.DeleteTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_DeleteTgChatId_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	repoMock.On("CheckUserExistence", mock.Anything, int64(123)).Return(true, nil)
	repoMock.On("DeleteChat", mock.Anything, int64(123)).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.DeleteTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_PostLinks_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	body := scrappertypes.AddLinkRequest{
		Link:    aws.String("https://github.com"),
		Tags:    &[]string{"tag1"},
		Filters: &[]string{"filter1"},
	}

	repoMock.On("SaveLink", mock.Anything, int64(123), mock.AnythingOfType("*domain.Link")).Return(nil)

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.PostLinks(c, scrappertypes.PostLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_PostLinks_InvalidLink(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	body := scrappertypes.AddLinkRequest{
		Link:    aws.String("test"),
		Tags:    &[]string{"tag1"},
		Filters: &[]string{"filter1"},
	}

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.PostLinks(c, scrappertypes.PostLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_PostLinks_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	body := scrappertypes.AddLinkRequest{
		Link:    aws.String("https://github.com"),
		Tags:    &[]string{"tag1"},
		Filters: &[]string{"filter1"},
	}

	repoMock.On("SaveLink", mock.Anything, int64(123), mock.AnythingOfType("*domain.Link")).Return(assert.AnError)

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.PostLinks(c, scrappertypes.PostLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_PostLinks_DuplicateLink(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	body := scrappertypes.AddLinkRequest{
		Link: aws.String("https://github.com"),
	}

	repoMock.On("SaveLink", mock.Anything, int64(123), mock.AnythingOfType("*domain.Link")).Return(nil).Twice()

	// First request.
	reqBody1, err := json.Marshal(body)
	assert.NoError(t, err)

	req1 := httptest.NewRequest(http.MethodPost, "/links?TgChatId=123", bytes.NewReader(reqBody1))
	req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec1 := httptest.NewRecorder()
	c1 := echo.New().NewContext(req1, rec1)
	err = h.PostLinks(c1, scrappertypes.PostLinksParams{TgChatId: 123})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// Second request.
	reqBody2, err := json.Marshal(body)
	assert.NoError(t, err)

	req2 := httptest.NewRequest(http.MethodPost, "/links?TgChatId=123", bytes.NewReader(reqBody2))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec2 := httptest.NewRecorder()
	c2 := echo.New().NewContext(req2, rec2)
	err = h.PostLinks(c2, scrappertypes.PostLinksParams{TgChatId: 123})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec2.Code)

	repoMock.AssertExpectations(t)
}

func Test_DeleteLinks_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	body := scrappertypes.RemoveLinkRequest{
		Link: aws.String("https://github.com"),
	}

	repoMock.On("DeleteLink", mock.Anything, int64(123), mock.AnythingOfType("*domain.Link")).Return(nil)

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.DeleteLinks(c, scrappertypes.DeleteLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_DeleteLinks_InvalidLink(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	body := scrappertypes.RemoveLinkRequest{
		Link: aws.String(""),
	}

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.DeleteLinks(c, scrappertypes.DeleteLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_DeleteLinks_LinkNotExist(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	body := scrappertypes.RemoveLinkRequest{
		Link: aws.String("test"),
	}

	repoMock.On("DeleteLink", mock.Anything, int64(123), mock.AnythingOfType("*domain.Link")).Return(&apperrors.LinkIsNotExistError{})

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.DeleteLinks(c, scrappertypes.DeleteLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_DeleteLinks_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	body := scrappertypes.RemoveLinkRequest{
		Link: aws.String("https://github.com"),
	}

	repoMock.On("DeleteLink", mock.Anything, int64(123), mock.AnythingOfType("*domain.Link")).Return(assert.AnError)

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.DeleteLinks(c, scrappertypes.DeleteLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_GetLinks_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	expectedLinks := []*domain.Link{
		{URL: "https://test", Tags: []string{"test_tag"}},
	}

	repoMock.On("GetListLinks", mock.Anything, int64(123)).Return(expectedLinks, nil)

	req := httptest.NewRequest(http.MethodGet, "/links?TgChatId=123", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.GetLinks(c, scrappertypes.GetLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp scrappertypes.ListLinksResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Len(t, *resp.Links, 1)
	assert.Equal(t, expectedLinks[0].URL, *(*resp.Links)[0].Url)
	assert.Equal(t, expectedLinks[0].Tags, *(*resp.Links)[0].Tags)
	repoMock.AssertExpectations(t)
}

func Test_GetLinks_EmptyList(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	repoMock.On("GetListLinks", mock.Anything, int64(123)).Return([]*domain.Link{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/links?TgChatId=123", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.GetLinks(c, scrappertypes.GetLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp scrappertypes.ListLinksResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Len(t, *resp.Links, 0)
	repoMock.AssertExpectations(t)
}

func Test_GetLinks_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := scrapperapi.NewScrapperHandler(repoMock, logger.NewDiscardLogger())

	repoMock.On("GetListLinks", mock.Anything, int64(123)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/links?TgChatId=123", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.GetLinks(c, scrappertypes.GetLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}
