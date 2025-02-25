package repository_test

import (
	"testing"
	"time"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
	"github.com/AFK068/bot/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
)

func TestRegisterChat(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	chatID := int64(1)

	t.Run("new chat registration", func(t *testing.T) {
		err := repo.RegisterChat(chatID)
		assert.NoError(t, err)
		assert.True(t, repo.CheckUserExistence(chatID))
	})

	t.Run("duplicate chat registration", func(t *testing.T) {
		err := repo.RegisterChat(chatID)
		assert.Error(t, err)
		assert.IsType(t, &apperrors.ChatAlreadyExistError{}, err)
	})
}

func TestSaveLink(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	chatID := int64(1)

	link := &domain.Link{URL: "https://github.com", UserAddID: chatID}

	t.Run("save link to new chat", func(t *testing.T) {
		err := repo.SaveLink(chatID, link)
		assert.NoError(t, err)

		links, err := repo.GetListLinks(chatID)
		assert.NoError(t, err)

		assert.Len(t, links, 1)
		assert.Equal(t, link.URL, links[0].URL)
	})

	t.Run("update existing link", func(t *testing.T) {
		newLink := &domain.Link{URL: "https://github.com", UserAddID: chatID, Tags: []string{"test"}}
		err := repo.SaveLink(chatID, newLink)
		assert.NoError(t, err)

		links, err := repo.GetListLinks(chatID)
		assert.NoError(t, err)

		assert.Len(t, links, 1)
		assert.Equal(t, []string{"test"}, links[0].Tags)
	})
}

func TestDeleteChat(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	chatID := int64(1)

	err := repo.RegisterChat(chatID)
	assert.NoError(t, err)

	t.Run("delete existing chat", func(t *testing.T) {
		err := repo.DeleteChat(chatID)
		assert.NoError(t, err)
		assert.False(t, repo.CheckUserExistence(chatID))
	})

	t.Run("delete non-existent chat", func(t *testing.T) {
		err := repo.DeleteChat(999)
		assert.NoError(t, err)
	})
}

func TestDeleteLink(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	chatID := int64(1)

	link := &domain.Link{URL: "https://github.com", UserAddID: chatID}

	err := repo.SaveLink(chatID, link)
	assert.NoError(t, err)

	t.Run("delete existing link", func(t *testing.T) {
		err := repo.DeleteLink(chatID, link)
		assert.NoError(t, err)

		links, _ := repo.GetListLinks(chatID)
		assert.Empty(t, links)
	})

	t.Run("delete non-existent link", func(t *testing.T) {
		err := repo.DeleteLink(chatID, &domain.Link{URL: "invalid"})
		assert.Error(t, err)
		assert.IsType(t, &apperrors.LinkIsNotExistError{}, err)
	})
}

func TestGetListLinks(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	chatID := int64(1)
	link1 := &domain.Link{URL: "https://github.com/1", UserAddID: chatID}
	link2 := &domain.Link{URL: "https://github.com/2", UserAddID: chatID}

	t.Run("empty list", func(t *testing.T) {
		links, err := repo.GetListLinks(chatID)
		assert.NoError(t, err)
		assert.Empty(t, links)
	})

	t.Run("non-empty list", func(t *testing.T) {
		err := repo.SaveLink(chatID, link1)
		assert.NoError(t, err)

		err = repo.SaveLink(chatID, link2)
		assert.NoError(t, err)

		links, err := repo.GetListLinks(chatID)
		assert.NoError(t, err)
		assert.Len(t, links, 2)
	})
}

func TestGetAllLinks(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()

	links := []*domain.Link{
		{URL: "https://github.com/1", UserAddID: 1},
		{URL: "https://github.com/2", UserAddID: 2},
		{URL: "https://stackoverflow.com/1", UserAddID: 1},
	}

	for _, l := range links {
		err := repo.SaveLink(l.UserAddID, l)
		assert.NoError(t, err)
	}

	allLinks := repo.GetAllLinks()
	assert.Len(t, allLinks, 3)

	for _, l := range links {
		assert.Contains(t, allLinks, l)
	}
}

func TestUpdateLastCheck(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	chatID := int64(1)
	link := &domain.Link{URL: "https://github.com", UserAddID: chatID}

	err := repo.SaveLink(chatID, link)
	assert.NoError(t, err)

	t.Run("successful update", func(t *testing.T) {
		beforeUpdate := time.Now().Add(-time.Second)
		err := repo.UpdateLastCheck(link)
		assert.NoError(t, err)

		updatedLink, _ := repo.GetListLinks(chatID)
		assert.True(t, updatedLink[0].LastCheck.After(beforeUpdate))
	})

	t.Run("update non-existent link", func(t *testing.T) {
		err := repo.UpdateLastCheck(&domain.Link{URL: "invalid"})
		assert.Error(t, err)
		assert.IsType(t, &apperrors.LinkIsNotExistError{}, err)
	})
}

func TestGetChatIDsByLink(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	link := &domain.Link{URL: "https://common.link", UserAddID: 1}

	err := repo.SaveLink(1, link)
	assert.NoError(t, err)

	err = repo.SaveLink(2, link)
	assert.NoError(t, err)

	err = repo.SaveLink(3, &domain.Link{URL: "https://unique.link", UserAddID: 3})
	assert.NoError(t, err)

	chatIDs := repo.GetChatIDsByLink(link)
	assert.ElementsMatch(t, []int64{1, 2}, chatIDs)
}
