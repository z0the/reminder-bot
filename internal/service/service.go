package service

import (
	"fmt"
	"log"
	"reminder-bot/internal/model"
	"reminder-bot/internal/service/keyboard"
	"reminder-bot/internal/service/keyboard/calendar"
	"reminder-bot/internal/service/keyboard/clock"
	"strings"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type RemindCreator interface {
}
type RemindList interface {
}
type Remind interface {
}
type Service struct {
	model *model.Model
	RemindCreator
	RemindList
	Remind
}

func NewService(mdl *model.Model) *Service {
	return &Service{model: mdl}
}
func CreateRemind(bot *botapi.BotAPI, updatesChan *botapi.UpdatesChannel, curChatID int64) {
	bot.Send(botapi.NewMessage(curChatID, "О чём вам напомнить"))
	type Remind struct {
		text string
		date time.Time
	}
	var curClock = make(map[string]int, 4)
	curClock["hour10"] = 1
	curClock["hour1"] = 2
	curClock["minute10"] = 0
	curClock["minute1"] = 0
	var curInlineMarkupID int
	var chosenDateMsgID int
	var remind Remind
	stepCount := 0
	curYear := time.Now().Year()
	curMonth := time.Now().Month()
	for update := range *updatesChan {
		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			bot.AnswerCallbackQuery(botapi.NewCallback(update.CallbackQuery.ID, data))
			log.Println("data: ", len(data))
			switch data {
			case "<":
				var keyboard botapi.InlineKeyboardMarkup
				keyboard, curYear, curMonth = calendar.HandlerPrevButton(curYear, curMonth)

				bot.Send(botapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))
			case ">":
				var keyboard botapi.InlineKeyboardMarkup
				keyboard, curYear, curMonth = calendar.HandlerNextButton(curYear, curMonth)

				bot.Send(botapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))

			}
			if strings.Contains(data, "+") {
				var keyboard botapi.InlineKeyboardMarkup
				keyboard, curClock = clock.PlusHandler(data, curClock)
				bot.Send(botapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))

			}
			if strings.Contains(data, "-") {
				var keyboard botapi.InlineKeyboardMarkup
				keyboard, curClock = clock.MinusHandler(data, curClock)
				bot.Send(botapi.NewEditMessageReplyMarkup(curChatID, curInlineMarkupID, keyboard))
			}
			if len(data) == 10 {
				// chosen = true
				newDate, err := time.Parse("2006.01.02", data)
				if err != nil {
					log.Println("Can't parse date from callback data, err:", err)
				}
				remind.date = newDate
				if chosenDateMsgID == 0 {
					msg := botapi.NewMessage(curChatID, fmt.Sprintf("Вы выбрали: %s\n", remind.date.Format("2006.01.02")))

					respMsg, err := bot.Send(msg)
					if err != nil {
						log.Println("err:", err)
					}
					chosenDateMsgID = respMsg.MessageID
				} else {
					log.Println("true")
					_, err := bot.Send(botapi.NewEditMessageText(curChatID, chosenDateMsgID, fmt.Sprintf("Вы выбрали: %s\n", remind.date.Format("2006.01.02"))))
					if err != nil {
						log.Println("err:", err)
					}
				}
			}
		}
		if update.Message != nil {
			if update.Message.Text == "Отмена" {
				msg := botapi.NewMessage(curChatID, "Создание напоминания отменено")
				msg.ReplyMarkup = keyboard.GetMainKeyboard()
				bot.Send(msg)
				return
			} else if update.Message.Text == "Далее" {
				stepCount++
			}
			switch stepCount {
			case 0:
				remind.text = update.Message.Text
				remind.date = time.Now()
				updateKeyboard := botapi.NewMessage(curChatID, "...")
				updateKeyboard.ReplyMarkup = botapi.NewReplyKeyboard([]botapi.KeyboardButton{botapi.NewKeyboardButton("Отмена"), botapi.NewKeyboardButton("Далее")})
				bot.Send(updateKeyboard)
				msg := botapi.NewMessage(update.Message.Chat.ID, "Когда нужно напомнить?")
				msg.ReplyMarkup = calendar.GenerateCalendar(curYear, curMonth)
				respMsg, err := bot.Send(msg)
				if err != nil {
					log.Println("Send err:", err)
				}
				curInlineMarkupID = respMsg.MessageID
			case 1:
				msg1 := botapi.NewMessage(curChatID, "Выберите время:")
				msg1.ReplyMarkup = clock.GenerateClockKeyboard(curClock["hour10"], curClock["hour1"], curClock["minute10"], curClock["minute1"])
				respMsg, err := bot.Send(msg1)
				if err != nil {
					log.Println("Send err:", err)
				}
				curInlineMarkupID = respMsg.MessageID
			case 2:
				hours := curClock["hour10"]*10 + curClock["hour1"]
				minutes := curClock["minute10"]*10 + curClock["minute1"]
				remind.date = remind.date.Add(time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes))
				// log.Println("Date: ", remind.date)
				msg1 := botapi.NewMessage(curChatID, fmt.Sprintf(`Ваше напоминание:

				%s

				%s в %s`, remind.text, remind.date.Format("2006.01.02"), remind.date.Format("15:04")))
				bot.Send(msg1)
				bot.Send(botapi.NewMessage(curChatID, "Если всё верно - нажмите Далее"))
			case 3:
				msg := botapi.NewMessage(curChatID, "Напоминание успешно создано!")
				msg.ReplyMarkup = keyboard.GetMainKeyboard()
				bot.Send(msg)
				return
			}
		}
	}
}
