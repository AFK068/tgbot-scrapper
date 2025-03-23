package bot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/AFK068/bot/internal/domain/apperrors"
	"github.com/AFK068/bot/pkg/utils"

	scrappertypes "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	command := msg.Command()

	b.Logger.Info("Received command", "chatID", chatID, "command", command)

	switch command {
	case StartCommand:
		b.handleStart(chatID)
	case HelpCommand:
		b.handleHelp(chatID)
	case TrackCommand:
		b.startTrackConversation(chatID)
	case UntrackCommand:
		b.handleUntrack(chatID, msg.CommandArguments())
	case ListCommand:
		b.handleList(chatID)
	default:
		b.SendMessage(chatID, "Unknown command. Use /help to see the list of available commands.")
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := msg.Text

	b.Logger.Info("Received message", "chatID", chatID, "text", text)

	conv := b.StateManager.GetConversation(chatID)
	if conv.FSM.Current() == ConversationStateIdle {
		b.SendMessage(chatID, "Please enter a command to start. Use /help to see the list of available commands.")
		return
	}

	switch conv.FSM.Current() {
	case ConversationStateAwaitingURL:
		if strings.Contains(text, "github.com") || strings.Contains(text, "stackoverflow.com") {
			conv.URL = text

			if err := conv.FSM.Event(context.Background(), EventSetURL); err != nil {
				b.Logger.Error("Error setting URL", "error", err)
				b.SendMessage(chatID, "Error setting URL. Please try again later.")

				return
			}

			b.SendMessage(chatID, "Enter tags separated by spaces (optional):", skipKeyboard)
		} else {
			b.SendMessage(chatID, "Invalid link. Please try again:")
		}

	case ConversationStateAwaitingTags:
		if text != "" && text != SkipOption {
			conv.Tags = strings.Split(text, " ")
		}

		if err := conv.FSM.Event(context.Background(), EventSetTags); err != nil {
			b.Logger.Error("Error setting tags", "error", err)
			b.SendMessage(chatID, "Error setting tags. Please try again later.")

			return
		}

		b.SendMessage(chatID, "Enter filters separated by spaces (optional):", skipKeyboard)

	case ConversationStateAwaitingFilter:
		if text != "" && text != SkipOption {
			conv.Filters = strings.Split(text, " ")
		}

		err := b.ScrapperClient.PostLinks(context.Background(), chatID, scrappertypes.AddLinkRequest{
			Link:    aws.String(conv.URL),
			Tags:    utils.SliceStringPtr(conv.Tags),
			Filters: utils.SliceStringPtr(conv.Filters),
		})

		if err != nil {
			b.Logger.Error("Error posting links", "error", err)
			b.handleError(chatID, err)
		} else {
			b.SendMessage(chatID, "Link successfully added!", mainKeyboard)
		}

		if err := conv.FSM.Event(context.Background(), EventComplete); err != nil {
			b.Logger.Error("Error completing tracking", "error", err)
			b.SendMessage(chatID, "Error completing tracking. Please try again later.", mainKeyboard)
		}

		b.StateManager.ClearConversation(chatID)
	}
}

func (b *Bot) startTrackConversation(chatID int64) {
	conv := b.StateManager.GetConversation(chatID)

	if err := conv.FSM.Event(context.Background(), EventStartTrack); err != nil {
		b.Logger.Error("Error starting tracking", "error", err)
		b.SendMessage(chatID, "Error starting tracking. Please try again later.")

		return
	}

	b.SendMessage(chatID, "Enter the link to track:", tgbotapi.NewRemoveKeyboard(true))
}

func (b *Bot) handleUntrack(chatID int64, link string) {
	if link == "" {
		b.SendMessage(chatID, "Specify the link to stop tracking: /untrack <link>", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	if err := b.ScrapperClient.DeleteLinks(context.Background(), chatID, scrappertypes.RemoveLinkRequest{
		Link: aws.String(link),
	}); err != nil {
		b.Logger.Error("Error deleting link", "error", err)
		b.handleError(chatID, err)
	} else {
		b.SendMessage(chatID, "Link successfully removed from tracking!", mainKeyboard)
	}
}

func (b *Bot) handleList(chatID int64) {
	links, err := b.ScrapperClient.GetLinks(context.Background(), chatID)
	if err != nil {
		b.Logger.Error("Error getting links", "error", err)
		b.handleError(chatID, err)

		return
	}

	if *links.Size == 0 {
		b.SendMessage(chatID, "No tracked links.")
		return
	}

	var builder strings.Builder

	builder.WriteString("Tracked links:\n")

	for _, link := range *links.Links {
		builder.WriteString(fmt.Sprintf("- %s\n", *link.Url))
	}

	b.SendMessage(chatID, builder.String())
}

func (b *Bot) handleStart(chatID int64) {
	if err := b.ScrapperClient.PostTgChatID(context.Background(), chatID); err != nil {
		b.Logger.Error("Error posting chat ID", "error", err)
		b.handleError(chatID, err)

		return
	}

	b.SendMessage(chatID, "Welcome! Use /help for a list of commands.", mainKeyboard)
}

func (b *Bot) handleHelp(chatID int64) {
	helpText := fmt.Sprintf(`Available commands:
/%s - %s
/%s - %s
/%s - %s
/%s - %s
/%s - %s`,
		StartCommand, StartCommandDescription,
		HelpCommand, HelpCommandDescription,
		TrackCommand, TrackCommandDescription,
		UntrackCommand, UntrackCommandDescription,
		ListCommand, ListCommandDescription,
	)

	b.SendMessage(chatID, helpText, mainKeyboard)
}

func (b *Bot) handleError(chatID int64, err error) {
	var errResp *apperrors.ErrorResponse
	if errors.As(err, &errResp) {
		switch errResp.Code {
		case http.StatusBadRequest:
			b.Logger.Error("Bad request error", "error", errResp.Message)
			b.SendMessage(chatID, fmt.Sprintf("❌ Request error: %s", errResp.Message))
		case http.StatusNotFound:
			b.Logger.Error("Not found error", "error", errResp.Message)
			b.SendMessage(chatID, fmt.Sprintf("🔍 Not found: %s", errResp.Message))
		case http.StatusUnauthorized:
			b.Logger.Error("Unauthorized access error", "error", errResp.Message)
			b.SendMessage(chatID, fmt.Sprintf("❌ Unauthorized access: %s", errResp.Message))
		default:
			b.Logger.Error("Unexpected error", "error", errResp.Message)
			b.SendMessage(chatID, "⚠️ An internal error occurred")
		}
	} else {
		b.Logger.Error("Internal error", "error", err)
		b.SendMessage(chatID, "⚠️ An internal error occurred")
	}
}
