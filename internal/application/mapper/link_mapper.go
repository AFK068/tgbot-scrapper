package mapper

import (
	"strings"
	"time"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"

	scrappertypes "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

func MapAddLinkRequestToDomain(tgChatID int64, addLinkRequest *scrappertypes.AddLinkRequest) (*domain.Link, error) {
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

	link.UserAddID = tgChatID

	switch {
	case strings.HasPrefix(*addLinkRequest.Link, "https://stackoverflow.com"):
		link.Type = domain.StackoverflowType
	case strings.HasPrefix(*addLinkRequest.Link, "https://github.com"):
		link.Type = domain.GithubType
	default:
		return nil, &apperrors.LinkTypeError{Message: "unsupported link type"}
	}

	link.LastCheck = time.Now()

	return link, nil
}
