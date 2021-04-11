package telegram

import (
	"reminder-bot/internal/database"
	"reminder-bot/internal/service/keyboard"
	"reminder-bot/internal/service/text"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *Bot) handlerListReminds(message *tgbotapi.Message) error {
	startMsg := tgbotapi.NewMessage(t.curChatID, "Список ваших напоминаний: ")
	startMsg.ReplyMarkup = keyboard.GetStopListViewKeyboard()
	t.bot.Send(startMsg)

	reminds, err := t.db.GetAllRemindesByChatID(t.curChatID)
	if err != nil {
		return err
	}
	for _, remind := range reminds {
		msg := tgbotapi.NewMessage(t.curChatID, text.RemindMessageText(&remind))
		msg.ReplyMarkup = keyboard.GetDeleteKeyboard(remind.ID)
		t.bot.Send(msg)
	}
	for update := range t.updatesChan {
		if update.CallbackQuery != nil {
			err := handleDeleteQuery(t.bot, t.db, &update, t.curChatID)
			if err != nil {
				return err
			}
		}
		if update.Message != nil {
			switch update.Message.Text {
			case "Прекратить просмотр":
				endMsg := tgbotapi.NewMessage(t.curChatID, "Просмотр ваших напоминаний окончен")
				endMsg.ReplyMarkup = keyboard.GetMainKeyboard()
				_, err := t.bot.Send(endMsg)
				return err
			default:
				t.handleUnknownCmd(update.Message)
			}
		}
	}
	return nil
}

func handleDeleteQuery(bot *tgbotapi.BotAPI, db database.BotDataBase, update *tgbotapi.Update, chatID int64) error {
	remindID, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	err = db.DeleteRemind(remindID)
	if err != nil {
		return err
	}
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "удалено"))
	bot.Send(tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID))
	// logrus.Info("message Id: ", update.CallbackQuery.Message.MessageID)
	return nil
}
