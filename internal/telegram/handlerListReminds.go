package telegram

import (
	"reminder-bot/internal/service/keyboard"
	"reminder-bot/internal/service/text"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *Bot) handlerListReminds(message *tgbotapi.Message) error {
	reminds, err := t.db.GetAllRemindesByChatID(t.curChatID)
	if err != nil {
		return err
	}
	for _, remind := range reminds {
		msg := tgbotapi.NewMessage(t.curChatID, text.RemindMessageText(&remind))
		msg.ReplyMarkup = keyboard.GetDeleteKeyboard()
		t.bot.Send(msg)
	}
	return nil
}
