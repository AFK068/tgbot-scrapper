package repository

import (
	"sync"
	"time"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
)

type InMemoryChatLinkRepository struct {
	links map[int64]map[string]*domain.Link
	mu    sync.RWMutex
}

func NewInMemoryLinkRepository() *InMemoryChatLinkRepository {
	return &InMemoryChatLinkRepository{
		links: make(map[int64]map[string]*domain.Link),
	}
}

func (r *InMemoryChatLinkRepository) RegisterChat(chatID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.links[chatID]; ok {
		return &apperrors.ChatAlreadyExistError{
			Message: "Chat is already exist",
		}
	}

	if _, ok := r.links[chatID]; !ok {
		r.links[chatID] = make(map[string]*domain.Link)
	}

	return nil
}

func (r *InMemoryChatLinkRepository) SaveLink(chatID int64, link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.links[chatID]; !ok {
		r.links[chatID] = make(map[string]*domain.Link)
	}

	r.links[chatID][link.URL] = link

	return nil
}

func (r *InMemoryChatLinkRepository) DeleteChat(chatID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.links, chatID)

	return nil
}

func (r *InMemoryChatLinkRepository) DeleteLink(chatID int64, link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.links[chatID][link.URL]; !ok {
		return &apperrors.LinkIsNotExistError{
			Message: "Link is not exist",
		}
	}

	delete(r.links[chatID], link.URL)

	return nil
}

func (r *InMemoryChatLinkRepository) GetListLinks(chatID int64) ([]*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	links := make([]*domain.Link, 0, len(r.links[chatID]))
	for _, link := range r.links[chatID] {
		links = append(links, link)
	}

	return links, nil
}

func (r *InMemoryChatLinkRepository) CheckUserExistence(chatID int64) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.links[chatID]; ok {
		return true
	}

	return false
}

func (r *InMemoryChatLinkRepository) GetAllLinks() []*domain.Link {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allLinks []*domain.Link

	for _, userLinks := range r.links {
		for _, link := range userLinks {
			allLinks = append(allLinks, link)
		}
	}

	return allLinks
}

func (r *InMemoryChatLinkRepository) UpdateLastCheck(link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.links[link.UserAddID][link.URL]; !ok {
		return &apperrors.LinkIsNotExistError{
			Message: "Link is not exist",
		}
	}

	r.links[link.UserAddID][link.URL].LastCheck = time.Now()

	return nil
}
