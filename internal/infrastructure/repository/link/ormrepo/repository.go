package ormrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
	"github.com/AFK068/bot/internal/infrastructure/repository/txs"
)

// Unfortunately, ORM in Go works very slowly and therefore it is not recommended to use the gorm.
// I understand that squirrel is not ORM, but the technical specifications said to use it.

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

	query, args, err := squirrel.Insert("tg_users").
		Columns("tg_id").
		Values(uid).
		Suffix("ON CONFLICT DO NOTHING RETURNING tg_id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	log.Print(query, args)

	var insertID int64

	err = querier.QueryRow(ctx, query, args...).Scan(&insertID)
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

	query, args, err := squirrel.Delete("tg_users").
		Where(squirrel.Eq{"tg_id": uid}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	tag, err := querier.Exec(ctx, query, args...)
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

	query, args, err := squirrel.Insert("links").
		Columns("url").
		Values(link.URL).
		Suffix("ON CONFLICT (url) DO NOTHING").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil
	}

	if _, err := querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("inserting link: %w", err)
	}

	var linkID int64

	query, args, err = squirrel.Select("id").
		From("links").
		Where(squirrel.Eq{"url": link.URL}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if err := querier.QueryRow(ctx, query, args...).Scan(&linkID); err != nil {
		return fmt.Errorf("getting link id: %w", err)
	}

	query, args, err = squirrel.Insert("user_link").
		Columns("tg_user_id", "link_id", "last_update", "filters", "tags").
		Values(uid, linkID, link.LastCheck, link.Filters, link.Tags).
		Suffix("ON CONFLICT (tg_user_id, link_id) DO UPDATE SET last_update = $3, filters = $4, tags = $5").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("inserting user link: %w", err)
	}

	return nil
}

func (r *Repository) DeleteLink(ctx context.Context, uid int64, link *domain.Link) error {
	querier := txs.GetQuerier(ctx, r.db)

	sub := squirrel.Select("id").
		From("links").
		Where(squirrel.Eq{"url": link.URL})

	query, args, err := squirrel.Delete("user_link").
		Where(squirrel.Eq{"tg_user_id": uid}).
		Where(squirrel.Expr("link_id = (?)", sub)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("building delete query: %w", err)
	}

	if _, err := querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("deleting link: %w", err)
	}

	return nil
}

func (r *Repository) GetListLinks(ctx context.Context, uid int64) ([]*domain.Link, error) {
	querier := txs.GetQuerier(ctx, r.db)

	query, args, err := squirrel.Select("l.url", "ul.last_update", "ul.filters", "ul.tags", "ul.tg_user_id").
		From("user_link ul").
		Join("links l ON ul.link_id = l.id").
		Where(squirrel.Eq{"ul.tg_user_id": uid}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := querier.Query(ctx, query, args...)
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

	subQuery := squirrel.Select("1").
		From("tg_users").
		Where(squirrel.Eq{"tg_id": uid})

	query, args, err := squirrel.Select().
		Column(squirrel.Expr("EXISTS (?)", subQuery)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return false, err
	}

	var exists bool
	if err := querier.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		return false, fmt.Errorf("checking user existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) GetChatIDsByLink(ctx context.Context, link *domain.Link) ([]int64, error) {
	querier := txs.GetQuerier(ctx, r.db)

	query, args, err := squirrel.Select("ul.tg_user_id").
		From("user_link ul").
		Join("links l ON ul.link_id = l.id").
		Where(squirrel.Eq{"l.url": link.URL}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := querier.Query(ctx, query, args...)
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

	// Update last check time to current time.
	newTime := r.TimeGetter()

	subQuery := squirrel.Select("id").
		From("links").
		Where(squirrel.Eq{"url": link.URL})

	query, args, err := squirrel.Update("user_link").
		Set("last_update", newTime).
		Where(squirrel.Eq{"tg_user_id": link.UserAddID}).
		Where(squirrel.Expr("link_id = (?)", subQuery)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	tag, err := querier.Exec(ctx, query, args...)
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

	query, args, err := squirrel.Select("l.url", "ul.last_update", "ul.filters", "ul.tags", "ul.tg_user_id").
		From("user_link ul").
		Join("links l ON ul.link_id = l.id").
		Where(squirrel.Eq{"tg_user_id": uid}).
		Where(squirrel.Expr("$2 = ANY(ul.tags)", tag)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := querier.Query(ctx, query, args...)
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
