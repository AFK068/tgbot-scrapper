package bot

import (
	"sync"
)

type StateManager struct {
	mu            sync.RWMutex
	conversations map[int64]*Conversation
}

func NewStateManager() *StateManager {
	return &StateManager{
		conversations: make(map[int64]*Conversation),
	}
}

func (sm *StateManager) GetConversation(chatID int64) *Conversation {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if conv, exists := sm.conversations[chatID]; exists {
		return conv
	}

	conv := NewConversation(chatID)
	sm.conversations[chatID] = conv

	return conv
}

func (sm *StateManager) ClearConversation(chatID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.conversations, chatID)
}
