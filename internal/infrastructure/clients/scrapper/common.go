package scrapper

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

func handleResponse(code int, body []byte) error {
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		var apiErr api.ApiErrorResponse
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}

		return fmt.Errorf("bad request: %s", *apiErr.Description)
	case http.StatusNotFound:
		var apiErr api.ApiErrorResponse
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}

		return fmt.Errorf("not found: %s", *apiErr.Description)
	default:
		return fmt.Errorf("unexpected status code: %d", code)
	}
}
