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
		return t.handleListReminds(message)
	case "Настройки":
		return t.handleSettings(message)
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
		if err != nil {
			return err
		}
		timeOffset := 0
		baseText := "Текущая разница в часах: "
		msg1 := tgbotapi.NewMessage(t.curChatID, fmt.Sprintf(baseText+"%d", timeOffset))
		msg1.ReplyMarkup = keyboard.GetSettingTimeOffsetKeyboard()
		respMsg, err := t.bot.Send(msg1)
		if err != nil {
			return err
		}
		editMsg := tgbotapi.NewEditMessageText(t.curChatID, respMsg.MessageID, "")
		editMsg.ReplyMarkup = keyboard.GetSettingTimeOffsetKeyboard()
		stepCount := 0
		for update := range t.updatesChan {
			if update.CallbackQuery != nil {
				data := update.CallbackQuery.Data
				_, err := t.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, data))
				if err != nil {
					return err
				}
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
					_, err := t.bot.Send(cancelMsg)
					if err != nil {
						return err
					}
					return nil
				} else if update.Message.Text == "Далее" {
					stepCount++
				}
				switch stepCount {
				case 1:
					var user = models.User{
						Name:           message.From.UserName,
						ChatID:         t.curChatID,
						TimeZoneOffset: 3 + timeOffset,
						AutoDelete:     false,
					}
					err := t.db.CreateUser(user)
					if err != nil {
						return err
					}
					msg1 := tgbotapi.NewMessage(t.curChatID, "Бот успешно настрен и готов к работе!")
					msg1.ReplyMarkup = keyboard.GetMainKeyboard()
					_, err = t.bot.Send(msg1)
					if err != nil {
						return err
					}
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
		if err != nil {
			return err
		}
	}
	return err
}
func (t *Bot) handleSettings(message *tgbotapi.Message) error {
	startMsg := tgbotapi.NewMessage(t.curChatID, "Ваши настройки:")
	startMsg.ReplyMarkup = keyboard.GetBackwardSaveKeyboard()
	_, err := t.bot.Send(startMsg)
	if err != nil {
		return err
	}
	user, err := t.db.GetUserByChatID(t.curChatID)
	if err != nil {
		return err
	}
	msg1 := tgbotapi.NewMessage(t.curChatID, fmt.Sprintf("Разница в часах с москвой: %d", user.TimeZoneOffset-3))
	msg1.ReplyMarkup = keyboard.GetSettingTimeOffsetKeyboard()
	respMsg, err := t.bot.Send(msg1)
	if err != nil {
		return err
	}
	var autoRem string
	if user.AutoDelete {
		autoRem = "включено"
	} else {
		autoRem = "отключено"
	}
	msg1.Text = fmt.Sprintf("Автоудаление отработавших напоминаний: %s", autoRem)
	var text string
	if user.AutoDelete {
		text = "отключить"
	} else {
		text = "включить"
	}
	row1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(text, "changeAutoDel"),
	}
	kb := tgbotapi.NewInlineKeyboardMarkup(row1)
	msg1.ReplyMarkup = kb
	respAD, err := t.bot.Send(msg1)
	if err != nil {
		return err
	}
	editMsg := tgbotapi.NewEditMessageText(t.curChatID, respMsg.MessageID, "")
	editMsg.ReplyMarkup = keyboard.GetSettingTimeOffsetKeyboard()
	editAD := tgbotapi.NewEditMessageText(t.curChatID, respAD.MessageID, "")
	timeOffset := user.TimeZoneOffset - 3
	baseText := "Разница в часах с москвой: "
	for update := range t.updatesChan {
		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			_, err := t.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, data))
			if err != nil {
				return err
			}
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
			case "changeAutoDel":
				user.AutoDelete = !user.AutoDelete

				if user.AutoDelete {
					autoRem = "включено"
				} else {
					autoRem = "отключено"
				}
				editAD.Text = fmt.Sprintf("Автоудаление отработавших напоминаний: %s", autoRem)
				if user.AutoDelete {
					text = "отключить"
				} else {
					text = "включить"
				}
				row1 := []tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData(text, "changeAutoDel"),
				}
				kb = tgbotapi.NewInlineKeyboardMarkup(row1)
				editAD.ReplyMarkup = &kb
				_, err := t.bot.Send(editAD)
				if err != nil {
					return err
				}
			}
		}
		if update.Message != nil {
			switch update.Message.Text {
			case "Назад":
				endMsg := tgbotapi.NewMessage(t.curChatID, "...")
				endMsg.ReplyMarkup = keyboard.GetMainKeyboard()
				_, err := t.bot.Send(endMsg)
				return err
			case "Сохранить":
				user.TimeZoneOffset = timeOffset + 3
				t.db.UpdateUser(user)
				t.db.UpdateUserBool(user, "auto_delete", user.AutoDelete)
				endMsg := tgbotapi.NewMessage(t.curChatID, "Настройки сохранены")
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

func (t *Bot) handleChangeQuery(update tgbotapi.Update, user models.User) error {
	_, err := t.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "изменение"))
	msg := tgbotapi.NewMessage(t.curChatID, "...")
	msg.ReplyMarkup = keyboard.GetCancelNextKeyboard()
	switch update.CallbackQuery.Data {
	case "changeTimeOffset":
		timeOffset := user.TimeZoneOffset - 3
		baseText := "Текущая разница в часах: "
		msg1 := tgbotapi.NewMessage(t.curChatID, fmt.Sprintf(baseText+"%d", timeOffset))
		msg1.ReplyMarkup = keyboard.GetSettingTimeOffsetKeyboard()
		respMsg, err := t.bot.Send(msg1)
		if err != nil {
			return err
		}
		editMsg := tgbotapi.NewEditMessageText(t.curChatID, respMsg.MessageID, "")
		editMsg.ReplyMarkup = keyboard.GetSettingTimeOffsetKeyboard()
		var stepCount int
		for update := range t.updatesChan {
			if update.CallbackQuery != nil {

			}
			if update.Message != nil {
				if update.Message.Text == "Отмена" {
					msg := tgbotapi.NewMessage(t.curChatID, "...")
					msg.ReplyMarkup = keyboard.GetMainKeyboard()
					_, err := t.bot.Send(msg)
					return err
				} else if update.Message.Text == "Далее" {
					stepCount++
				}
				if stepCount == 1 {
					endMsg := tgbotapi.NewMessage(t.curChatID, "Разница в часах с москвой изменена")
					endMsg.ReplyMarkup = keyboard.GetMainKeyboard()
					_, err := t.bot.Send(endMsg)
					return err
				}
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (t *Bot) handleUnknownCmd(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(t.curChatID, "")

	msg.Text = "Не знаю таких команд :("
	_, err := t.bot.Send(msg)
	return err
}
