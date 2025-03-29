package ormrepo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/AFK068/bot/internal/config"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
	"github.com/AFK068/bot/internal/infrastructure/repository/link/ormrepo"
	"github.com/AFK068/bot/internal/testcontainer"
)

const (
	TestConfigPath = "../../../../../config/test.yaml"
)

func setupDB(t *testing.T) (*ormrepo.Repository, *pgxpool.Pool, context.Context) {
	ctx := context.Background()

	config, err := config.NewConfig(TestConfigPath)
	assert.NoError(t, err)

	testContainer, err := testcontainer.NewPostgresTestcontainerContainer(ctx, config)
	assert.NoError(t, err)

	dbPool, cleanup, err := testContainer.SetupTestPostgresContainer(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cleanup())
	})

	repo := ormrepo.NewRepository(dbPool)

	return repo, dbPool, ctx
}

func Test_RegisterChat_Success(t *testing.T) {
	repo, dbPool, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	var count int
	err = dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM tg_users WHERE tg_id = $1", uid).Scan(&count)

	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func Test_RegisterChat_ChatAlreadyExist_Failure(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	err = repo.RegisterChat(ctx, uid)
	assert.Error(t, err)
	assert.IsType(t, &apperrors.ChatAlreadyExistError{}, err)
}

func Test_DeleteChat_Success(t *testing.T) {
	repo, dbPool, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	err = repo.DeleteChat(ctx, uid)
	assert.NoError(t, err)

	var count int
	err = dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM tg_users WHERE tg_id = $1", uid).Scan(&count)

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func Test_DeleteChat_ChatNotFound_Failure(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.DeleteChat(ctx, uid)
	assert.Error(t, err)
}

func Test_SaveLink_Success(t *testing.T) {
	repo, dbPool, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	link := &domain.Link{
		URL:       "https://github.com/AFK068/bot",
		LastCheck: time.Now(),
		Filters:   []string{"filter1", "filter2"},
		Tags:      []string{"go", "bot"},
	}

	err = repo.SaveLink(ctx, uid, link)
	assert.NoError(t, err)

	var count int
	err = dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM user_link WHERE tg_user_id = $1", uid).Scan(&count)

	assert.NoError(t, err)
	assert.Equal(t, count, 1)
}

func Test_DeleteLink_Success(t *testing.T) {
	repo, dbPool, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	link := &domain.Link{
		URL:       "https://github.com/AFK068/bot",
		LastCheck: time.Now(),
		Filters:   []string{"filter1", "filter2"},
		Tags:      []string{"go", "bot"},
	}

	err = repo.SaveLink(ctx, uid, link)
	assert.NoError(t, err)

	err = repo.DeleteLink(ctx, uid, link)
	assert.NoError(t, err)

	var count int
	err = dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM user_link WHERE tg_user_id = $1", uid).Scan(&count)

	assert.NoError(t, err)
	assert.Equal(t, count, 0)
}

func Test_GetListLinks_Success(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	timeGetter := func() time.Time {
		return time.Time{}
	}

	link := &domain.Link{
		URL:       "https://github.com/AFK068/bot",
		LastCheck: timeGetter(),
		Filters:   []string{"filter1", "filter2"},
		Tags:      []string{"go", "bot"},
	}

	err = repo.SaveLink(ctx, uid, link)
	assert.NoError(t, err)

	links, err := repo.GetListLinks(ctx, uid)
	assert.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, link.URL, links[0].URL)
	assert.Equal(t, link.LastCheck, links[0].LastCheck)
	assert.Equal(t, link.Filters, links[0].Filters)
	assert.Equal(t, link.Tags, links[0].Tags)
}

func Test_GetListLinks_NoLinksFound_Success(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	links, err := repo.GetListLinks(ctx, uid)
	assert.NoError(t, err)
	assert.Len(t, links, 0)
}

func Test_CheckUserExistence_Success(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	exist, err := repo.CheckUserExistence(ctx, uid)
	assert.NoError(t, err)
	assert.True(t, exist)
}

func Test_CheckUserExistence_UserNotFound_Failure(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	exist, err := repo.CheckUserExistence(ctx, uid)
	assert.NoError(t, err)
	assert.False(t, exist)
}

func Test_GetChatIDsByLink_Success(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	timeGetter := func() time.Time {
		return time.Time{}
	}

	link := &domain.Link{
		URL:       "https://github.com/AFK068/bot",
		LastCheck: timeGetter(),
		Filters:   []string{"filter1", "filter2"},
		Tags:      []string{"go", "bot"},
	}

	err = repo.SaveLink(ctx, uid, link)
	assert.NoError(t, err)

	chatIDs, err := repo.GetChatIDsByLink(ctx, link)
	assert.NoError(t, err)

	assert.Len(t, chatIDs, 1)
	assert.Equal(t, uid, chatIDs[0])
}

func Test_UpdateLastCheck_Success(t *testing.T) {
	repo, dbPool, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	testTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	repo.TimeGetter = func() time.Time {
		return testTime
	}

	link := &domain.Link{
		UserAddID: uid,
		URL:       "https://github.com/AFK068/bot",
		Filters:   []string{"filter1", "filter2"},
		Tags:      []string{"go", "bot"},
	}

	err = repo.SaveLink(ctx, uid, link)
	assert.NoError(t, err)

	err = repo.UpdateLastCheck(ctx, link)
	assert.NoError(t, err)

	q := `
	SELECT last_update
	FROM user_link
	WHERE tg_user_id = $1 AND link_id = (SELECT id FROM links WHERE url = $2);
	`

	var lastCheck time.Time
	err = dbPool.QueryRow(ctx, q, uid, link.URL).Scan(&lastCheck)

	assert.NoError(t, err)
	assert.Equal(t, testTime, lastCheck)
}

func Test_GetLinksByTag_Success(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	testTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	repo.TimeGetter = func() time.Time {
		return testTime
	}

	link := &domain.Link{
		UserAddID: uid,
		URL:       "https://github.com/AFK068/bot",

		Filters: []string{"filter1", "filter2"},
		Tags:    []string{"go", "bot"},
	}

	err = repo.SaveLink(ctx, uid, link)
	assert.NoError(t, err)

	links, err := repo.GetLinksByTag(ctx, uid, "go")
	assert.NoError(t, err)

	assert.Len(t, links, 1)
	assert.Equal(t, link.URL, links[0].URL)
	assert.Equal(t, link.LastCheck, links[0].LastCheck)
	assert.Equal(t, link.Filters, links[0].Filters)
	assert.Equal(t, link.Tags, links[0].Tags)
}

func TestGetLinksPagination_Success(t *testing.T) {
	repo, _, ctx := setupDB(t)

	uid := int64(12345)

	err := repo.RegisterChat(ctx, uid)
	assert.NoError(t, err)

	countLinks := 10
	limit := 2

	for i := 0; i < countLinks; i++ {
		link := &domain.Link{
			UserAddID: uid,
			URL:       fmt.Sprintf("%d", i),
		}

		err = repo.SaveLink(ctx, uid, link)
		assert.NoError(t, err)
	}

	offset := 0

	for offset < countLinks {
		pagedLinks, err := repo.GetLinksPagination(ctx, uint64(offset), uint64(limit)) //nolint
		assert.NoError(t, err)

		if offset+limit > countLinks {
			assert.Len(t, pagedLinks, countLinks-offset)
		} else {
			assert.Len(t, pagedLinks, limit)
		}

		for i, link := range pagedLinks {
			expectedURL := fmt.Sprintf("%d", offset+i)

			assert.Equal(t, uid, link.UserAddID)
			assert.Equal(t, expectedURL, link.URL)
		}

		offset += limit
	}
}
