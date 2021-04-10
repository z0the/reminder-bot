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
	for update := range t.updatesChan {
		if update.CallbackQuery != nil {
			err := handleDeleteQuery(t.bot, &update)
			if err != nil {
				return err
			}
		}
		if update.Message != nil {

		}
	}
	return nil
}

func handleDeleteQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	return nil
}
