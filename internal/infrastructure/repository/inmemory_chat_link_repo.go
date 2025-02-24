package repository

import (
	"sync"
	"time"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
)

type InMemoryChatLinkRepository struct {
	Links map[int64]map[string]*domain.Link
	mu    sync.RWMutex
}

func NewInMemoryLinkRepository() *InMemoryChatLinkRepository {
	return &InMemoryChatLinkRepository{
		Links: make(map[int64]map[string]*domain.Link),
	}
}

func (r *InMemoryChatLinkRepository) RegisterChat(chatID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Links[chatID]; ok {
		return &apperrors.ChatAlreadyExistError{
			Message: "Chat is already exist",
		}
	}

	if _, ok := r.Links[chatID]; !ok {
		r.Links[chatID] = make(map[string]*domain.Link)
	}

	return nil
}

func (r *InMemoryChatLinkRepository) SaveLink(chatID int64, link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Links[chatID]; !ok {
		r.Links[chatID] = make(map[string]*domain.Link)
	}

	r.Links[chatID][link.URL] = link

	return nil
}

func (r *InMemoryChatLinkRepository) DeleteChat(chatID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Links, chatID)

	return nil
}

func (r *InMemoryChatLinkRepository) DeleteLink(chatID int64, link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Links[chatID][link.URL]; !ok {
		return &apperrors.LinkIsNotExistError{
			Message: "Link is not exist",
		}
	}

	delete(r.Links[chatID], link.URL)

	return nil
}

func (r *InMemoryChatLinkRepository) GetListLinks(chatID int64) ([]*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	links := make([]*domain.Link, 0, len(r.Links[chatID]))
	for _, link := range r.Links[chatID] {
		links = append(links, link)
	}

	return links, nil
}

func (r *InMemoryChatLinkRepository) CheckUserExistence(chatID int64) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.Links[chatID]; ok {
		return true
	}

	return false
}

func (r *InMemoryChatLinkRepository) GetAllLinks() []*domain.Link {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allLinks []*domain.Link

	for _, userLinks := range r.Links {
		for _, link := range userLinks {
			allLinks = append(allLinks, link)
		}
	}

	return allLinks
}

func (r *InMemoryChatLinkRepository) UpdateLastCheck(link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Links[link.UserAddID][link.URL]; !ok {
		return &apperrors.LinkIsNotExistError{
			Message: "Link is not exist",
		}
	}

	r.Links[link.UserAddID][link.URL].LastCheck = time.Now()

	return nil
}

func (r *InMemoryChatLinkRepository) GetChatIDsByLink(link *domain.Link) []int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var chatIDs []int64

	for chatID, userLinks := range r.Links {
		if _, ok := userLinks[link.URL]; ok {
			chatIDs = append(chatIDs, chatID)
		}
	}

	return chatIDs
}
