package scrapper

import (
	"encoding/json"
	"net/http"

	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	"github.com/AFK068/bot/internal/domain/apperrors"
)

func handleResponse(code int, body []byte) error {
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest, http.StatusNotFound, http.StatusUnauthorized:
		var apiErr api.ApiErrorResponse
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return &apperrors.ErrorResponse{
				Code:    code,
				Message: "failed to decode error response",
			}
		}

		return &apperrors.ErrorResponse{
			Code:    code,
			Message: *apiErr.Description,
		}
	default:
		return &apperrors.ErrorResponse{
			Code:    code,
			Message: "unexpected error",
		}
	}
}
