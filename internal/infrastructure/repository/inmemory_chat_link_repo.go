package repository

import (
	"context"
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

func (r *InMemoryChatLinkRepository) RegisterChat(_ context.Context, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Links[uid]; ok {
		return &apperrors.ChatAlreadyExistError{
			Message: "Chat is already exist",
		}
	}

	if _, ok := r.Links[uid]; !ok {
		r.Links[uid] = make(map[string]*domain.Link)
	}

	return nil
}

func (r *InMemoryChatLinkRepository) SaveLink(_ context.Context, uid int64, link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Links[uid]; !ok {
		r.Links[uid] = make(map[string]*domain.Link)
	}

	r.Links[uid][link.URL] = link

	return nil
}

func (r *InMemoryChatLinkRepository) DeleteChat(_ context.Context, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Links, uid)

	return nil
}

func (r *InMemoryChatLinkRepository) DeleteLink(_ context.Context, uid int64, link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Links[uid][link.URL]; !ok {
		return &apperrors.LinkIsNotExistError{
			Message: "Link is not exist",
		}
	}

	delete(r.Links[uid], link.URL)

	return nil
}

func (r *InMemoryChatLinkRepository) GetListLinks(_ context.Context, uid int64) ([]*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	links := make([]*domain.Link, 0, len(r.Links[uid]))
	for _, link := range r.Links[uid] {
		links = append(links, link)
	}

	return links, nil
}

func (r *InMemoryChatLinkRepository) CheckUserExistence(_ context.Context, uid int64) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.Links[uid]; ok {
		return true, nil
	}

	return false, nil
}

func (r *InMemoryChatLinkRepository) GetAllLinks(_ context.Context) ([]*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allLinks []*domain.Link

	for _, userLinks := range r.Links {
		for _, link := range userLinks {
			allLinks = append(allLinks, link)
		}
	}

	return allLinks, nil
}

func (r *InMemoryChatLinkRepository) UpdateLastCheck(_ context.Context, link *domain.Link) error {
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

func (r *InMemoryChatLinkRepository) GetChatIDsByLink(_ context.Context, link *domain.Link) ([]int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var chatIDs []int64

	for chatID, userLinks := range r.Links {
		if _, ok := userLinks[link.URL]; ok {
			chatIDs = append(chatIDs, chatID)
		}
	}

	return chatIDs, nil
}
