package mapper

import (
	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
)

func MapAddLinkRequestToDomain(addLinkRequest *api.AddLinkRequest) (*domain.Link, error) {
	if addLinkRequest.Link == nil || *addLinkRequest.Link == "" {
		return nil, &apperrors.LinkValidateError{Message: "link is required"}
	}

	link := &domain.Link{
		URL: *addLinkRequest.Link,
	}

	if addLinkRequest.Tags != nil {
		link.Tags = *addLinkRequest.Tags
	}

	if addLinkRequest.Filters != nil {
		link.Filters = *addLinkRequest.Filters
	}

	switch *addLinkRequest.Link {
	case "https://api.stackexchange.com":
		link.Type = domain.StackoverflowType
	case "https://api.github.com":
		link.Type = domain.GithubType
	default:
		return nil, &apperrors.LinkTypeError{Message: "unsupported link type"}
	}

	return link, nil
}
