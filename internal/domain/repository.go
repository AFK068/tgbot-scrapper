package domain

import "context"

type ChatLinkRepository interface {
	// Chat methods.
	RegisterChat(ctx context.Context, uid int64) error
	DeleteChat(ctx context.Context, uid int64) error

	// Link methods.
	SaveLink(ctx context.Context, uid int64, link *Link) error
	DeleteLink(ctx context.Context, uid int64, link *Link) error
	GetListLinks(ctx context.Context, uid int64) ([]*Link, error)
	CheckUserExistence(ctx context.Context, uid int64) (bool, error)
	GetChatIDsByLink(ctx context.Context, link *Link) ([]int64, error)
	UpdateLastCheck(ctx context.Context, link *Link) error
	GetLinksByTag(ctx context.Context, uid int64, tag string) ([]*Link, error)
	GetLinksPagination(ctx context.Context, offset, limit uint64) ([]*Link, error)
}
