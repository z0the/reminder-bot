package telegram

import (
	"fmt"
	"reminder-bot/internal/models"
	"reminder-bot/internal/service/keyboard"
	"reminder-bot/internal/service/keyboard/calendar"
	"reminder-bot/internal/service/keyboard/clock"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	curClock          = make(map[string]int, 4)
	curChatID         int64
	curInlineMarkupID int
	chosenDateMsgID   int
	remind            = new(models.Remind)
	curYear           int
	curMonth          time.Month
	stepCount         = 0
)

func (t *Bot) handleCreateRemind(message *tgbotapi.Message) error {
	curClock["hour10"] = 1
	curClock["hour1"] = 2
	curClock["minute10"] = 0
	curClock["minute1"] = 0
	curChatID = t.curChatID
	curYear = time.Now().Year()
	curMonth = time.Now().Month()
	_, err := t.bot.Send(tgbotapi.NewMessage(curChatID, "О чём вам напомнить"))
	if err != nil {
		return err
	}
	for update := range t.updatesChan {
		if update.CallbackQuery != nil {
			err := handleInlineQueryCreating(t.bot, &update)
			if err != nil {
				return err
			}
		}
		if update.Message != nil {
			if update.Message.Text == "Отмена" {
				msg := tgbotapi.NewMessage(curChatID, "Создание напоминания отменено")
				msg.ReplyMarkup = keyboard.GetMainKeyboard()
				_, err := t.bot.Send(msg)
				return err
			} else if update.Message.Text == "Далее" {
				stepCount++
			}
			err := t.processMessage(&update)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func handleInlineQueryCreating(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	data := update.CallbackQuery.Data
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, data))
	var keyboard tgbotapi.InlineKeyboardMarkup
	var inlineKeyboard tgbotapi.InlineKeyboardMarkup

	switch data {
	case "<":
		keyboard, curYear, curMonth = calendar.HandlerPrevButton(curYear, curMonth)
		_, err := bot.Send(tgbotapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))
		return err
	case ">":
		keyboard, curYear, curMonth = calendar.HandlerNextButton(curYear, curMonth)
		_, err := bot.Send(tgbotapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))
		return err
	}
	if strings.Contains(data, "+") {
		keyboard, curClock = clock.PlusHandler(data, curClock)
		_, err := bot.Send(tgbotapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))
		return err
	}
	if strings.Contains(data, "-") {
		inlineKeyboard, curClock = clock.MinusHandler(data, curClock)
		_, err := bot.Send(tgbotapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, inlineKeyboard))
		return err
	}
	if len(data) == 10 {
		return updateMessageWithDate(bot, data)
	}
	return nil
}

func updateMessageWithDate(bot *tgbotapi.BotAPI, date string) error {
	newDate, err := time.Parse("2006.01.02", date)
	if err != nil {
		return err
	}
	remind.ActivationTime = time.Date(newDate.Year(), newDate.Month(), newDate.Day(), 0, 0, 0, 0, time.Now().UTC().Location())

	if chosenDateMsgID == 0 {
		msg := tgbotapi.NewMessage(curChatID, fmt.Sprintf("Вы выбрали: %s\n", remind.ActivationTime.Format("2006.01.02")))

		respMsg, err := bot.Send(msg)
		if err != nil {
			return err
		}
		chosenDateMsgID = respMsg.MessageID
	} else {
		_, err := bot.Send(tgbotapi.NewEditMessageText(curChatID, chosenDateMsgID, fmt.Sprintf("Вы выбрали: %s\n", remind.ActivationTime.Format("2006.01.02"))))
		return err
	}
	return nil
}

func (t *Bot) processMessage(update *tgbotapi.Update) error {
	switch stepCount {
	case 0:
		remind.Text = update.Message.Text
		remind.ActivationTime = time.Time{}
		remind.ChatID = curChatID
		lastID, err := t.db.GetLastRemindIDByChatID(curChatID)
		if err != nil {
			return err
		}
		remind.IDForChat = lastID + 1
		updateKeyboard := tgbotapi.NewMessage(curChatID, "...")
		updateKeyboard.ReplyMarkup = keyboard.GetCancelNextKeyboard()
		t.bot.Send(updateKeyboard)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Когда нужно напомнить?")
		msg.ReplyMarkup = calendar.GenerateCalendar(curYear, curMonth)
		respMsg, err := t.bot.Send(msg)
		if err != nil {
			return err
		}
		curInlineMarkupID = respMsg.MessageID
	case 1:
		msg1 := tgbotapi.NewMessage(curChatID, "Выберите время:")
		msg1.ReplyMarkup = clock.GenerateClockKeyboard(curClock["hour10"], curClock["hour1"], curClock["minute10"], curClock["minute1"])
		respMsg, err := t.bot.Send(msg1)
		if err != nil {
			return err
		}
		curInlineMarkupID = respMsg.MessageID
	case 2:
		hours := curClock["hour10"]*10 + curClock["hour1"]
		minutes := curClock["minute10"]*10 + curClock["minute1"]
		remind.ActivationTime = remind.ActivationTime.Add(time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes))
		user, err := t.db.GetUserByChatID(curChatID)
		if err != nil {
			return err
		}
		remind.ActivationTime = remind.ActivationTime.Add(-time.Duration(user.TimeZoneOffset) * time.Hour)
		msg1 := tgbotapi.NewMessage(curChatID, fmt.Sprintf(`Ваше напоминание:

				%s

				%s в %s`, remind.Text, remind.ActivationTime.Format("2006.01.02"), remind.ActivationTime.Format("15:04")))
		t.bot.Send(msg1)
		t.bot.Send(tgbotapi.NewMessage(curChatID, "Если всё верно - нажмите Далее"))
	case 3:
		remindRes, err := t.db.CreateRemind(*remind)
		if err != nil {
			return err
		}
		if time.Until(remind.ActivationTime) < time.Hour {
			go t.serveRemind(remindRes)
		}
		msg := tgbotapi.NewMessage(curChatID, "Напоминание успешно создано!")
		msg.ReplyMarkup = keyboard.GetMainKeyboard()
		_, err = t.bot.Send(msg)
		return err
	}
	return nil
}
