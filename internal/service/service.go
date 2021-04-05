package service

import (
	"fmt"
	"reminder-bot/internal/entities"
	"reminder-bot/internal/model"
	"reminder-bot/internal/service/keyboard"
	"reminder-bot/internal/service/keyboard/calendar"
	"reminder-bot/internal/service/keyboard/clock"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type RemindCreator interface {
	CreateRemind(remind entities.Remind) (int, error)
}
type RemindList interface {
}

// type Remind interface {
// }
type Service struct {
	RemindCreator
	RemindList
	// Remind
}

func NewService(mdl *model.Model) *Service {
	return &Service{
		RemindCreator: NewCreateService(mdl),
	}
}
func CreateRemind(bot *tgbotapi.BotAPI, updatesChan *tgbotapi.UpdatesChannel, curChatID int64) {
	bot.Send(tgbotapi.NewMessage(curChatID, "О чём вам напомнить"))

	var curClock = make(map[string]int, 4)
	curClock["hour10"] = 1
	curClock["hour1"] = 2
	curClock["minute10"] = 0
	curClock["minute1"] = 0
	var curInlineMarkupID int
	var chosenDateMsgID int
	var remind = new(entities.Remind)
	stepCount := 0
	curYear := time.Now().Year()
	curMonth := time.Now().Month()
	for update := range *updatesChan {
		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, data))
			logrus.Println("data: ", len(data))
			switch data {
			case "<":
				var keyboard tgbotapi.InlineKeyboardMarkup
				keyboard, curYear, curMonth = calendar.HandlerPrevButton(curYear, curMonth)

				bot.Send(tgbotapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))
			case ">":
				var keyboard tgbotapi.InlineKeyboardMarkup
				keyboard, curYear, curMonth = calendar.HandlerNextButton(curYear, curMonth)

				bot.Send(tgbotapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))

			}
			if strings.Contains(data, "+") {
				var keyboard tgbotapi.InlineKeyboardMarkup
				keyboard, curClock = clock.PlusHandler(data, curClock)
				bot.Send(tgbotapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))

			}
			if strings.Contains(data, "-") {
				var keyboard tgbotapi.InlineKeyboardMarkup
				keyboard, curClock = clock.MinusHandler(data, curClock)
				bot.Send(tgbotapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))
			}
			if len(data) == 10 {
				// chosen = true
				newDate, err := time.Parse("2006.01.02", data)
				if err != nil {
					logrus.Println("Can't parse date from callback data, err:", err)
				}
				remind.ActivationTime = newDate
				if chosenDateMsgID == 0 {
					msg := tgbotapi.NewMessage(curChatID, fmt.Sprintf("Вы выбрали: %s\n", remind.ActivationTime.Format("2006.01.02")))

					respMsg, err := bot.Send(msg)
					if err != nil {
						logrus.Println("err:", err)
					}
					chosenDateMsgID = respMsg.MessageID
				} else {
					logrus.Println("true")
					_, err := bot.Send(tgbotapi.NewEditMessageText(curChatID, chosenDateMsgID, fmt.Sprintf("Вы выбрали: %s\n", remind.ActivationTime.Format("2006.01.02"))))
					if err != nil {
						logrus.Println("err:", err)
					}
				}
			}
		}
		if update.Message != nil {
			if update.Message.Text == "Отмена" {
				msg := tgbotapi.NewMessage(curChatID, "Создание напоминания отменено")
				msg.ReplyMarkup = keyboard.GetMainKeyboard()
				bot.Send(msg)
				return
			} else if update.Message.Text == "Далее" {
				stepCount++
			}
			switch stepCount {
			case 0:
				remind.Text = update.Message.Text
				remind.ActivationTime = time.Now()
				updateKeyboard := tgbotapi.NewMessage(curChatID, "...")
				updateKeyboard.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("Отмена"), tgbotapi.NewKeyboardButton("Далее")})
				bot.Send(updateKeyboard)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Когда нужно напомнить?")
				msg.ReplyMarkup = calendar.GenerateCalendar(curYear, curMonth)
				respMsg, err := bot.Send(msg)
				if err != nil {
					logrus.Println("Send err:", err)
				}
				curInlineMarkupID = respMsg.MessageID
			case 1:
				msg1 := tgbotapi.NewMessage(curChatID, "Выберите время:")
				msg1.ReplyMarkup = clock.GenerateClockKeyboard(curClock["hour10"], curClock["hour1"], curClock["minute10"], curClock["minute1"])
				respMsg, err := bot.Send(msg1)
				if err != nil {
					logrus.Println("Send err:", err)
				}
				curInlineMarkupID = respMsg.MessageID
			case 2:
				hours := curClock["hour10"]*10 + curClock["hour1"]
				minutes := curClock["minute10"]*10 + curClock["minute1"]
				remind.ActivationTime = remind.ActivationTime.Add(time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes))
				// logrus.Println("Date: ", remind.ActivationTime)
				msg1 := tgbotapi.NewMessage(curChatID, fmt.Sprintf(`Ваше напоминание:

				%s

				%s в %s`, remind.Text, remind.ActivationTime.Format("2006.01.02"), remind.ActivationTime.Format("15:04")))
				bot.Send(msg1)
				bot.Send(tgbotapi.NewMessage(curChatID, "Если всё верно - нажмите Далее"))
			case 3:
				msg := tgbotapi.NewMessage(curChatID, "Напоминание успешно создано!")
				msg.ReplyMarkup = keyboard.GetMainKeyboard()
				bot.Send(msg)
				return
			}
		}
	}
}
