package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

const (
	cmdStart = "start"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case cmdStart:
		return b.handleStartCmd(message)
	default:
		return b.handleUnknownCmd(message)
	}
}
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Привет"))
}

func (b *Bot) handleStartCmd(message *tgbotapi.Message) error {
	logrus.Info("Bot is started...")
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	msg.Text = `
		<b>Привет юзер :)</b>
		Используй меня, чтобы создавать для себя напоминания о самых важных вещах).
		`
	msg.ParseMode = "HTML"

	_, err := b.bot.Send(msg)
	return err
}
func (b *Bot) handleUnknownCmd(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	msg.Text = "Не знаю таких команд :("
	_, err := b.bot.Send(msg)
	return err
}
