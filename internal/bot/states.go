package bot

import (
	"context"

	"github.com/looplab/fsm"
)

const (
	StateIdle           = "idle"
	StateAwaitingURL    = "awaiting_url"
	StateAwaitingTags   = "awaiting_tags"
	StateAwaitingFilter = "awaiting_filter"
)

type Conversation struct {
	ChatID  int64
	URL     string
	Tags    []string
	Filters []string
	FSM     *fsm.FSM
}

func NewConversation(chatID int64) *Conversation {
	conversation := &Conversation{
		ChatID: chatID,
	}

	conversation.FSM = fsm.NewFSM(
		StateIdle,
		fsm.Events{
			{Name: "start_track", Src: []string{StateIdle}, Dst: StateAwaitingURL},
			{Name: "set_url", Src: []string{StateAwaitingURL}, Dst: StateAwaitingTags},
			{Name: "set_tags", Src: []string{StateAwaitingTags}, Dst: StateAwaitingFilter},
			{Name: "complete", Src: []string{StateAwaitingFilter}, Dst: StateIdle},
		},
		fsm.Callbacks{
			"enter_state": func(_ context.Context, _ *fsm.Event) {},
		},
	)

	return conversation
}
