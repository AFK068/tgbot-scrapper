package bot

import (
	"context"

	"github.com/looplab/fsm"
)

const (
	ConversationStateIdle           = "idle"
	ConversationStateAwaitingURL    = "awaiting_url"
	ConversationStateAwaitingTags   = "awaiting_tags"
	ConversationStateAwaitingFilter = "awaiting_filter"

	EventStartTrack = "start_track"
	EventSetURL     = "set_url"
	EventSetTags    = "set_tags"
	EventComplete   = "complete"

	EnterState = "enter_state"
)

type Conversation struct {
	ChatID  int64
	URL     string
	Tags    []string
	Filters []string
	FSM     *fsm.FSM
}

func NewConversationWithFSM(chatID int64) *Conversation {
	conversation := &Conversation{
		ChatID: chatID,
	}

	conversation.FSM = fsm.NewFSM(
		ConversationStateIdle,
		fsm.Events{
			{Name: EventStartTrack, Src: []string{ConversationStateIdle}, Dst: ConversationStateAwaitingURL},
			{Name: EventSetURL, Src: []string{ConversationStateAwaitingURL}, Dst: ConversationStateAwaitingTags},
			{Name: EventSetTags, Src: []string{ConversationStateAwaitingTags}, Dst: ConversationStateAwaitingFilter},
			{Name: EventComplete, Src: []string{ConversationStateAwaitingFilter}, Dst: ConversationStateIdle},
		},
		fsm.Callbacks{
			EnterState: func(_ context.Context, _ *fsm.Event) {},
		},
	)

	return conversation
}
