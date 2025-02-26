package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	repomock "github.com/AFK068/bot/internal/domain/mocks"
	"github.com/AFK068/bot/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthLinkMiddleware_Success(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	mw := middleware.AuthLinkMiddleware(repoMock)

	req := httptest.NewRequest(http.MethodGet, "/links", http.NoBody)
	rec := httptest.NewRecorder()

	req.Header.Set("Tg-Chat-Id", "1")

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/links")

	repoMock.On("CheckUserExistence", int64(1)).Return(true)

	called := false
	nextHandler := func(c echo.Context) error {
		called = true
		return c.String(http.StatusOK, "OK")
	}

	err := mw(nextHandler)(c)

	assert.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestAuthLinkMiddleware_SkipNonLinksPath(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	mw := middleware.AuthLinkMiddleware(repoMock)

	req := httptest.NewRequest(http.MethodGet, "/other", http.NoBody)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)

	called := false
	nextHandler := func(c echo.Context) error {
		called = true
		return c.String(http.StatusOK, "OK")
	}

	err := mw(nextHandler)(c)

	assert.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestAuthLinkMiddleware_MissingHeader(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	mw := middleware.AuthLinkMiddleware(repoMock)

	req := httptest.NewRequest(http.MethodGet, "/links", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/links")

	err := mw(func(_ echo.Context) error { return nil })(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestAuthLinkMiddleware_InvalidHeader(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	mw := middleware.AuthLinkMiddleware(repoMock)

	req := httptest.NewRequest(http.MethodGet, "/links", http.NoBody)
	rec := httptest.NewRecorder()

	req.Header.Set("Tg-Chat-Id", "no_int")

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/links")

	err := mw(func(_ echo.Context) error { return nil })(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	repoMock.AssertExpectations(t)
}

func TestAuthLinkMiddleware_UserNotExist(t *testing.T) {
	repoMock := repomock.NewChatLinkRepository(t)
	mw := middleware.AuthLinkMiddleware(repoMock)

	req := httptest.NewRequest(http.MethodGet, "/links", http.NoBody)
	rec := httptest.NewRecorder()

	req.Header.Set("Tg-Chat-Id", "123")

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/links")

	repoMock.On("CheckUserExistence", int64(123)).Return(false)

	err := mw(func(_ echo.Context) error { return nil })(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	repoMock.AssertExpectations(t)
}
