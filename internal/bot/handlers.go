package bot

import (
	"context"
	"fmt"
	"strings"

	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	"github.com/AFK068/bot/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StartCommand            = "start"
	StartCommandDescription = "Start command"

	HelpCommand            = "help"
	HelpCommandDescription = "List available commands"

	TrackCommand            = "track"
	TrackCommandDescription = "Start tracking a link"

	UntrackCommand            = "untrack"
	UntrackCommandDescription = "Stop tracking a link"

	ListCommand            = "list"
	ListCommandDescription = "Show list of tracked links"
)

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	command := msg.Command()

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

	conv := b.StateManager.GetConversation(chatID)
	if conv.FSM.Current() == StateIdle {
		b.SendMessage(chatID, "Please enter a command to start. Use /help to see the list of available commands.")
		return
	}

	switch conv.FSM.Current() {
	case StateAwaitingURL:
		if strings.Contains(text, "github.com") || strings.Contains(text, "stackoverflow.com") {
			conv.URL = text

			if err := conv.FSM.Event(context.Background(), "set_url"); err != nil {
				b.SendMessage(chatID, "Error setting URL. Please try again later.")
				return
			}

			b.SendMessage(chatID, "Enter tags separated by spaces (optional):")
		} else {
			b.SendMessage(chatID, "Invalid link. Please try again:")
		}

	case StateAwaitingTags:
		if text != "" {
			conv.Tags = strings.Split(text, " ")
		}

		if err := conv.FSM.Event(context.Background(), "set_tags"); err != nil {
			b.SendMessage(chatID, "Error setting tags. Please try again later.")
			return
		}

		b.SendMessage(chatID, "Enter filters separated by spaces (optional):")

	case StateAwaitingFilter:
		if text != "" {
			conv.Filters = strings.Split(text, " ")
		}

		err := b.ScrapperClient.PostLinks(context.Background(), chatID, api.AddLinkRequest{
			Link:    aws.String(conv.URL),
			Tags:    utils.SliceStringPtr(conv.Tags),
			Filters: utils.SliceStringPtr(conv.Filters),
		})

		if err != nil {
			b.SendMessage(chatID, "Error adding link. Please try again later.")
		} else {
			b.SendMessage(chatID, "Link successfully added!")
		}

		if err := conv.FSM.Event(context.Background(), "complete"); err != nil {
			b.SendMessage(chatID, "Error completing tracking. Please try again later.")
		}

		b.StateManager.ClearConversation(chatID)
	}
}

func (b *Bot) startTrackConversation(chatID int64) {
	conv := b.StateManager.GetConversation(chatID)

	if err := conv.FSM.Event(context.Background(), "start_track"); err != nil {
		b.SendMessage(chatID, "Error starting tracking. Please try again later.")
		return
	}

	b.SendMessage(chatID, "Enter the link to track:")
}

func (b *Bot) handleUntrack(chatID int64, link string) {
	if link == "" {
		b.SendMessage(chatID, "Specify the link to stop tracking: /untrack <link>")
		return
	}

	if err := b.ScrapperClient.DeleteLinks(context.Background(), chatID, api.RemoveLinkRequest{
		Link: aws.String(link),
	}); err != nil {
		b.SendMessage(chatID, "Error removing link. Please try again later.")
	} else {
		b.SendMessage(chatID, "Link successfully removed from tracking!")
	}
}

func (b *Bot) handleList(chatID int64) {
	links, err := b.ScrapperClient.GetLinks(context.Background(), chatID)
	if err != nil {
		b.SendMessage(chatID, "Error retrieving list of links.")
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
		b.SendMessage(chatID, "Registration error. Please try again later.")
		return
	}

	b.SendMessage(chatID, "Welcome! Use /help for a list of commands.")
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

	b.SendMessage(chatID, helpText)
}
