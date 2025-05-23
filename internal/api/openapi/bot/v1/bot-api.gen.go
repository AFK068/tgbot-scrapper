// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package v1

import (
	"time"

	"github.com/labstack/echo/v4"
)

// Defines values for LinkUpdateType.
const (
	GithubIssue           LinkUpdateType = "github_issue"
	GithubPullRequest     LinkUpdateType = "github_pull_request"
	GithubRepository      LinkUpdateType = "github_repository"
	StackoverflowAnswer   LinkUpdateType = "stackoverflow_answer"
	StackoverflowComment  LinkUpdateType = "stackoverflow_comment"
	StackoverflowQuestion LinkUpdateType = "stackoverflow_question"
)

// ApiErrorResponse defines model for ApiErrorResponse.
type ApiErrorResponse struct {
	Code             *string   `json:"code,omitempty"`
	Description      *string   `json:"description,omitempty"`
	ExceptionMessage *string   `json:"exceptionMessage,omitempty"`
	ExceptionName    *string   `json:"exceptionName,omitempty"`
	Stacktrace       *[]string `json:"stacktrace,omitempty"`
}

// LinkUpdate defines model for LinkUpdate.
type LinkUpdate struct {
	Type        *LinkUpdateType `json:"Type,omitempty"`
	UserName    *string         `json:"UserName,omitempty"`
	Description *string         `json:"description,omitempty"`
	Id          *int64          `json:"id,omitempty"`
	TgChatIds   *[]int64        `json:"tgChatIds,omitempty"`
	Url         *string         `json:"url,omitempty"`
	СreatedAt   *time.Time      `json:"сreatedAt,omitempty"`
}

// LinkUpdateType defines model for LinkUpdate.Type.
type LinkUpdateType string

// PostUpdatesJSONRequestBody defines body for PostUpdates for application/json ContentType.
type PostUpdatesJSONRequestBody = LinkUpdate

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Отправить обновление
	// (POST /updates)
	PostUpdates(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostUpdates converts echo context to params.
func (w *ServerInterfaceWrapper) PostUpdates(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostUpdates(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/updates", wrapper.PostUpdates)

}
