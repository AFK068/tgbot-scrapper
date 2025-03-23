package scrapper

import (
	"encoding/json"
	"net/http"

	"github.com/AFK068/bot/internal/domain/apperrors"

	scrappertypes "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

func (c *Client) handleResponse(code int, body []byte) error {
	switch code {
	case http.StatusOK:
		c.Logger.Info("Request successful")
		return nil
	case http.StatusBadRequest, http.StatusNotFound, http.StatusUnauthorized:
		var apiErr scrappertypes.ApiErrorResponse
		if err := json.Unmarshal(body, &apiErr); err != nil {
			c.Logger.Error("Failed to decode error response", "error", err)

			return &apperrors.ErrorResponse{
				Code:    code,
				Message: "failed to decode error response",
			}
		}

		c.Logger.Warn("API error response", "code", code, "description", *apiErr.Description)

		return &apperrors.ErrorResponse{
			Code:    code,
			Message: *apiErr.Description,
		}
	default:
		c.Logger.Error("Unexpected error", "code", code)

		return &apperrors.ErrorResponse{
			Code:    code,
			Message: "unexpected error",
		}
	}
}
