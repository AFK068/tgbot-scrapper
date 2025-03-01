package bot

import (
	"fmt"

	"github.com/AFK068/bot/config"
	"github.com/AFK068/bot/internal/infrastructure/clients/scrapper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service interface {
	Run() error
	SendMessage(chatID int64, text string, replyMarkup ...interface{})
}

type Bot struct {
	API            *tgbotapi.BotAPI
	Config         *config.BotConfig
	ScrapperClient *scrapper.Client
	StateManager   *StateManager
}

func NewBot(cfg *config.BotConfig, sc *scrapper.Client) *Bot {
	return &Bot{
		Config:         cfg,
		ScrapperClient: sc,
	}
}

func (b *Bot) Run() error {
	if err := b.setBotAPI(); err != nil {
		return fmt.Errorf("setting bot api: %w", err)
	}

	if err := b.setBotCommands(); err != nil {
		return fmt.Errorf("setting bot commands: %w", err)
	}

	b.StateManager = NewStateManager()

	updates := b.initUpdatesChannel()
	go b.processUpdates(updates)

	return nil
}

func (b *Bot) SendMessage(chatID int64, text string, replyMarkup ...interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)

	if len(replyMarkup) > 0 {
		if keyboard, ok := replyMarkup[0].(tgbotapi.ReplyKeyboardMarkup); ok {
			msg.ReplyMarkup = keyboard
		}

		if keyboard, ok := replyMarkup[0].(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = keyboard
		}

		if _, ok := replyMarkup[0].(tgbotapi.ReplyKeyboardRemove); ok {
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
	}

	if _, err := b.API.Send(msg); err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
	}
}

func (b *Bot) processUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
		} else {
			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) setBotAPI() error {
	botAPI, err := tgbotapi.NewBotAPI(b.Config.Token)
	if err != nil {
		return fmt.Errorf("creating bot api: %w", err)
	}

	b.API = botAPI

	return nil
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.API.GetUpdatesChan(u)
}

func (b *Bot) initBotCommands() tgbotapi.SetMyCommandsConfig {
	commands := []tgbotapi.BotCommand{
		{
			Command:     StartCommand,
			Description: StartCommandDescription,
		},
		{
			Command:     HelpCommand,
			Description: HelpCommandDescription,
		},
		{
			Command:     TrackCommand,
			Description: TrackCommandDescription,
		},
		{
			Command:     UntrackCommand,
			Description: UntrackCommandDescription,
		},
		{
			Command:     ListCommand,
			Description: ListCommandDescription,
		},
	}

	return tgbotapi.SetMyCommandsConfig{
		Commands: commands,
	}
}

func (b *Bot) setBotCommands() error {
	commandsConfig := b.initBotCommands()

	_, err := b.API.Request(commandsConfig)
	if err != nil {
		return fmt.Errorf("setting bot commands: %w", err)
	}

	return nil
}
