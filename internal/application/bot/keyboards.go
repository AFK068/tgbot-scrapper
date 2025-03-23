package bot

import (
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

	SkipOption = "Skip"
)

var (
	mainKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/"+TrackCommand),
			tgbotapi.NewKeyboardButton("/"+ListCommand),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/"+HelpCommand),
			tgbotapi.NewKeyboardButton("/"+UntrackCommand),
		),
	)

	skipKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(SkipOption),
		),
	)
)
