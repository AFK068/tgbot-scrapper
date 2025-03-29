package sqlrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
	"github.com/AFK068/bot/internal/infrastructure/repository/txs"
)

type timeGetter func() time.Time

type Repository struct {
	TimeGetter timeGetter
	db         *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:         db,
		TimeGetter: time.Now,
	}
}

func (r *Repository) RegisterChat(ctx context.Context, uid int64) error {
	querier := txs.GetQuerier(ctx, r.db)

	query := `INSERT INTO tg_users (tg_id) VALUES ($1) ON CONFLICT DO NOTHING RETURNING tg_id;`

	var insertID int64

	err := querier.QueryRow(ctx, query, uid).Scan(&insertID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &apperrors.ChatAlreadyExistError{Message: "Chat is already exist"}
		}

		return err
	}

	return nil
}

func (r *Repository) DeleteChat(ctx context.Context, uid int64) error {
	querier := txs.GetQuerier(ctx, r.db)

	query := `DELETE FROM tg_users WHERE tg_id = $1;`

	tag, err := querier.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return &apperrors.ChatIsNotExistError{Message: "Chat not found"}
	}

	return err
}

func (r *Repository) SaveLink(ctx context.Context, uid int64, link *domain.Link) error {
	querier := txs.GetQuerier(ctx, r.db)

	query := `INSERT INTO links (url) VALUES ($1) ON CONFLICT (url) DO NOTHING;`
	if _, err := querier.Exec(ctx, query, link.URL); err != nil {
		return fmt.Errorf("inserting link: %w", err)
	}

	var linkID int64

	query = `SELECT id FROM links WHERE url = $1;`
	if err := querier.QueryRow(ctx, query, link.URL).Scan(&linkID); err != nil {
		return fmt.Errorf("getting link id: %w", err)
	}

	query = `
	INSERT INTO user_link (tg_user_id, link_id, last_update, filters, tags)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (tg_user_id, link_id) DO UPDATE
	SET last_update = $3, filters = $4, tags = $5;
	`

	if _, err := querier.Exec(ctx, query, uid, linkID, link.LastCheck, link.Filters, link.Tags); err != nil {
		return fmt.Errorf("inserting user link: %w", err)
	}

	return nil
}

func (r *Repository) DeleteLink(ctx context.Context, uid int64, link *domain.Link) error {
	querier := txs.GetQuerier(ctx, r.db)

	query := `
	DELETE FROM user_link 
	WHERE tg_user_id = $1 AND link_id = (SELECT id FROM links WHERE url = $2);
	`

	if _, err := querier.Exec(ctx, query, uid, link.URL); err != nil {
		return fmt.Errorf("deleting link: %w", err)
	}

	return nil
}

func (r *Repository) GetListLinks(ctx context.Context, uid int64) ([]*domain.Link, error) {
	querier := txs.GetQuerier(ctx, r.db)

	query := `
	SELECT l.url, ul.last_update, ul.filters, ul.tags, ul.tg_user_id
	FROM user_link ul
	JOIN links l ON ul.link_id = l.id
	WHERE ul.tg_user_id = $1;
	`

	rows, err := querier.Query(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("getting links: %w", err)
	}

	defer rows.Close()

	var links []*domain.Link

	for rows.Next() {
		var link domain.Link

		if err := rows.Scan(&link.URL, &link.LastCheck, &link.Filters, &link.Tags, &link.UserAddID); err != nil {
			return nil, fmt.Errorf("scanning link: %w", err)
		}

		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over rows: %w", err)
	}

	return links, nil
}

func (r *Repository) CheckUserExistence(ctx context.Context, uid int64) (bool, error) {
	querier := txs.GetQuerier(ctx, r.db)

	query := `SELECT EXISTS (SELECT 1 FROM tg_users WHERE tg_id = $1);`

	var exists bool
	if err := querier.QueryRow(ctx, query, uid).Scan(&exists); err != nil {
		return false, fmt.Errorf("checking user existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) GetChatIDsByLink(ctx context.Context, link *domain.Link) ([]int64, error) {
	querier := txs.GetQuerier(ctx, r.db)

	query := `
	SELECT ul.tg_user_id
	FROM user_link ul
	INNER JOIN links l ON ul.link_id = l.id
	WHERE l.url = $1;
	`

	rows, err := querier.Query(ctx, query, link.URL)
	if err != nil {
		return nil, fmt.Errorf("getting chat ids by link: %w", err)
	}

	defer rows.Close()

	var chatIDs []int64

	for rows.Next() {
		var chatID int64

		if err := rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("scanning chat id: %w", err)
		}

		chatIDs = append(chatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over rows: %w", err)
	}

	return chatIDs, nil
}

func (r *Repository) UpdateLastCheck(ctx context.Context, link *domain.Link) error {
	querier := txs.GetQuerier(ctx, r.db)

	query := `
	UPDATE user_link 
	SET last_update = $1 
	WHERE tg_user_id = $2 AND link_id = (SELECT id FROM links WHERE url = $3);
	`

	// Update last check time to current time.
	newTime := r.TimeGetter()

	tag, err := querier.Exec(ctx, query, newTime, link.UserAddID, link.URL)
	if err != nil {
		return fmt.Errorf("updating last check: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return &apperrors.LinkIsNotExistError{Message: "Link is not exist"}
	}

	return nil
}

func (r *Repository) GetLinksByTag(ctx context.Context, uid int64, tag string) ([]*domain.Link, error) {
	querier := txs.GetQuerier(ctx, r.db)

	query := ` 
	SELECT l.url, ul.last_update, ul.filters, ul.tags, ul.tg_user_id
	FROM user_link ul
	JOIN links l ON ul.link_id = l.id
	WHERE tg_user_id = $1 AND $2 = ANY(ul.tags);
	`

	rows, err := querier.Query(ctx, query, uid, tag)
	if err != nil {
		return nil, fmt.Errorf("getting links by tag: %w", err)
	}

	defer rows.Close()

	var links []*domain.Link

	for rows.Next() {
		var link domain.Link

		if err := rows.Scan(&link.URL, &link.LastCheck, &link.Filters, &link.Tags, &link.UserAddID); err != nil {
			return nil, fmt.Errorf("scanning link: %w", err)
		}

		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over rows: %w", err)
	}

	return links, nil
}

func (r *Repository) GetLinksPagination(ctx context.Context, offset, limit uint64) ([]*domain.Link, error) {
	querier := txs.GetQuerier(ctx, r.db)

	query := `
	SELECT l.url, ul.last_update, ul.filters, ul.tags, ul.tg_user_id
	FROM user_link ul
	JOIN links l ON ul.link_id = l.id
	LIMIT $1 OFFSET $2;
	`

	rows, err := querier.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("getting links pagination: %w", err)
	}

	defer rows.Close()

	var links []*domain.Link

	for rows.Next() {
		var link domain.Link

		if err := rows.Scan(&link.URL, &link.LastCheck, &link.Filters, &link.Tags, &link.UserAddID); err != nil {
			return nil, fmt.Errorf("scanning link: %w", err)
		}

		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over rows: %w", err)
	}

	return links, nil
}
