package telegram

import (
	"fmt"
	"reminder-bot/internal/models"
	"reminder-bot/internal/service/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

const (
	cmdStart = "start"
)

func (t *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case cmdStart:
		return t.handleStartCmd(message)
	default:
		return t.handleUnknownCmd(message)
	}
}
func (t *Bot) handleMessage(message *tgbotapi.Message) error {
	switch message.Text {
	case "Создать напоминание":
		return t.handleCreateRemind(message)
	case "Посмотреть список моих напоминаний":
		return t.handlerListReminds(message)
	default:
		return t.handleUnknownCmd(message)
	}
}

func (t *Bot) handleStartCmd(message *tgbotapi.Message) error {
	logrus.Info("Bot is started...")
	msg := tgbotapi.NewMessage(t.curChatID, "")
	msg.ParseMode = "HTML"
	users, err := t.db.GetAllUsers()
	if err != nil {
		return err
	}
	var registered bool
	for _, user := range users {
		if user.ChatID == t.curChatID {
			registered = true
		}
	}
	if !registered {
		msg.Text = `
		<b>Привет юзер :)</b>

		Используй меня, чтобы создавать для себя напоминания о самых важных вещах).

		<b>Для начала работы бота необходимо указать твою разницу во времени с Москвой.</b>
		`
		msg.ReplyMarkup = keyboard.GetCancelNextKeyboard()
		_, err = t.bot.Send(msg)

		timeOffset := 0
		baseText := "Текущая разница в часах: "
		msg1 := tgbotapi.NewMessage(t.curChatID, fmt.Sprintf(baseText+"%d", timeOffset))
		msg1.ReplyMarkup = keyboard.GetSettingTimeOffsetKeyboard()
		respMsg, _ := t.bot.Send(msg1)
		editMsg := tgbotapi.NewEditMessageText(t.curChatID, respMsg.MessageID, "")
		editMsg.ReplyMarkup = keyboard.GetSettingTimeOffsetKeyboard()
		stepCount := 0
		for update := range t.updatesChan {
			if update.CallbackQuery != nil {
				data := update.CallbackQuery.Data
				t.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, data))
				// logrus.Println("data: ", data)
				switch data {
				case "+":
					timeOffset++
					if timeOffset > 0 {
						editMsg.Text = fmt.Sprintf(baseText+"+%d", timeOffset)
					} else {
						editMsg.Text = fmt.Sprintf(baseText+"%d", timeOffset)
					}
					_, err := t.bot.Send(editMsg)
					if err != nil {
						return err
					}
				case "-":
					timeOffset--
					if timeOffset > 0 {
						editMsg.Text = fmt.Sprintf(baseText+"+%d", timeOffset)
					} else {
						editMsg.Text = fmt.Sprintf(baseText+"%d", timeOffset)
					}
					_, err := t.bot.Send(editMsg)
					if err != nil {
						return err
					}
				}
			}
			if update.Message != nil {
				if update.Message.Text == "Отмена" {
					cancelMsg := tgbotapi.NewMessage(t.curChatID, "Запуск бота отменён :(")
					t.bot.Send(cancelMsg)
					return nil
				} else if update.Message.Text == "Далее" {
					stepCount++
				}
				switch stepCount {
				case 1:
					var user = models.User{
						Name:           message.From.UserName,
						ChatID:         t.curChatID,
						TimeZoneOffset: 3+timeOffset,
					}
					err := t.db.CreateUser(user)
					if err != nil {
						return err
					}
					msg1 := tgbotapi.NewMessage(t.curChatID, "Бот успешно настрен и готов к работе!")
					msg1.ReplyMarkup = keyboard.GetMainKeyboard()
					t.bot.Send(msg1)
					return nil
				}
			}
		}

	} else {
		msg.Text = `
		<b>Привет юзер!</b>
		Для тебя Бот уже запущен :)
		`
		msg.ReplyMarkup = keyboard.GetMainKeyboard()
		_, err = t.bot.Send(msg)
	}
	return err
}
func (t *Bot) handleUnknownCmd(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(t.curChatID, "")

	msg.Text = "Не знаю таких команд :("
	_, err := t.bot.Send(msg)
	return err
}
