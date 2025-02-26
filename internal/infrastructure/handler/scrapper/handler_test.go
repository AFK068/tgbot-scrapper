package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	scrapperapi "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
	repomock "github.com/AFK068/bot/internal/domain/mocks"
	handler "github.com/AFK068/bot/internal/infrastructure/handler/scrapper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostTgChatId_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	repoMock.On("CheckUserExistence", int64(123)).Return(false)
	repoMock.On("RegisterChat", int64(123)).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.PostTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestPostTgChatId_AlreadyExists(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	repoMock.On("CheckUserExistence", int64(123)).Return(true)

	req := httptest.NewRequest(http.MethodPost, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.PostTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestPostTgChatId_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	repoMock.On("CheckUserExistence", int64(123)).Return(false)
	repoMock.On("RegisterChat", int64(123)).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.PostTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteTgChatId_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	repoMock.On("CheckUserExistence", int64(123)).Return(true)
	repoMock.On("DeleteChat", int64(123)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.DeleteTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteTgChatId_UserNotFound(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	repoMock.On("CheckUserExistence", int64(123)).Return(false)

	req := httptest.NewRequest(http.MethodDelete, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.DeleteTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteTgChatId_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	repoMock.On("CheckUserExistence", int64(123)).Return(true)
	repoMock.On("DeleteChat", int64(123)).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/tg-chat/123", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err := h.DeleteTgChatId(c, 123)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestPostLinks_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	body := scrapperapi.AddLinkRequest{
		Link:    aws.String("https://github.com"),
		Tags:    &[]string{"tag1"},
		Filters: &[]string{"filter1"},
	}

	repoMock.On("SaveLink", int64(123), mock.AnythingOfType("*domain.Link")).Return(nil)

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.PostLinks(c, scrapperapi.PostLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestPostLinks_InvalidLink(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	body := scrapperapi.AddLinkRequest{
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

	err = h.PostLinks(c, scrapperapi.PostLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestPostLinks_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	body := scrapperapi.AddLinkRequest{
		Link:    aws.String("https://github.com"),
		Tags:    &[]string{"tag1"},
		Filters: &[]string{"filter1"},
	}

	repoMock.On("SaveLink", int64(123), mock.AnythingOfType("*domain.Link")).Return(assert.AnError)

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.PostLinks(c, scrapperapi.PostLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteLinks_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	body := scrapperapi.RemoveLinkRequest{
		Link: aws.String("https://github.com"),
	}

	repoMock.On("DeleteLink", int64(123), mock.AnythingOfType("*domain.Link")).Return(nil)

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.DeleteLinks(c, scrapperapi.DeleteLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteLinks_InvalidLink(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	body := scrapperapi.RemoveLinkRequest{
		Link: aws.String(""),
	}

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.DeleteLinks(c, scrapperapi.DeleteLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteLinks_LinkNotExist(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	body := scrapperapi.RemoveLinkRequest{
		Link: aws.String("test"),
	}

	repoMock.On("DeleteLink", int64(123), mock.AnythingOfType("*domain.Link")).Return(&apperrors.LinkIsNotExistError{})

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.DeleteLinks(c, scrapperapi.DeleteLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteLinks_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	body := scrapperapi.RemoveLinkRequest{
		Link: aws.String("https://github.com"),
	}

	repoMock.On("DeleteLink", int64(123), mock.AnythingOfType("*domain.Link")).Return(assert.AnError)

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/links?TgChatId=123", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.DeleteLinks(c, scrapperapi.DeleteLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestGetLinks_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	expectedLinks := []*domain.Link{
		{URL: "https://test", Tags: []string{"test_tag"}},
	}

	repoMock.On("GetListLinks", int64(123)).Return(expectedLinks, nil)

	req := httptest.NewRequest(http.MethodGet, "/links?TgChatId=123", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.GetLinks(c, scrapperapi.GetLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp scrapperapi.ListLinksResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Len(t, *resp.Links, 1)
	assert.Equal(t, expectedLinks[0].URL, *(*resp.Links)[0].Url)
	assert.Equal(t, expectedLinks[0].Tags, *(*resp.Links)[0].Tags)
	repoMock.AssertExpectations(t)
}

func TestGetLinks_EmptyList(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	repoMock.On("GetListLinks", int64(123)).Return([]*domain.Link{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/links?TgChatId=123", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.GetLinks(c, scrapperapi.GetLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp scrapperapi.ListLinksResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Len(t, *resp.Links, 0)
	repoMock.AssertExpectations(t)
}

func TestGetLinks_Failure(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	h := handler.NewScrapperHandler(repoMock)

	repoMock.On("GetListLinks", int64(123)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/links?TgChatId=123", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.GetLinks(c, scrapperapi.GetLinksParams{TgChatId: 123})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}
