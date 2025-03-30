package bot

import (
	"context"

	"github.com/looplab/fsm"
)

// Why alias?
// Because looplab/fsm only works with strings and
// I think it's bad idea every time to cast to a strings.

type ConversationState = string

type Event = string

const (
	ConversationStateIdle           ConversationState = "idle"
	ConversationStateAwaitingURL    ConversationState = "awaiting_url"
	ConversationStateAwaitingTags   ConversationState = "awaiting_tags"
	ConversationStateAwaitingFilter ConversationState = "awaiting_filter"

	EventStartTrack Event = "start_track"
	EventSetURL     Event = "set_url"
	EventSetTags    Event = "set_tags"
	EventComplete   Event = "complete"

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
