package handler

import (
	"reminder-bot/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	bot          *tgbotapi.BotAPI
	mainKeyboard *tgbotapi.ReplyKeyboardMarkup
	updatesChan  *tgbotapi.UpdatesChannel
	curUpdate    *tgbotapi.Update
	curChatID    int64
	service      *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(bot *tgbotapi.BotAPI, mainKeyboard *tgbotapi.ReplyKeyboardMarkup, updatesChan *tgbotapi.UpdatesChannel) *Handler {
	handler := new(Handler)
	handler.updatesChan = updatesChan
	handler.bot = bot
	handler.mainKeyboard = mainKeyboard
	return handler
}

func (h *Handler) Handle(update *tgbotapi.Update) {
	h.curUpdate = update
	h.curChatID = update.Message.Chat.ID

	switch update.Message.Text {
	case "/start":
		logrus.Println("Bot is started...")
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
			<b>Привет юзер :)</b>
			Используй меня, чтобы создавать для себя напоминания о самых важных вещах).`)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = h.mainKeyboard

		h.bot.Send(msg)

	case "Создать напоминание":
		service.CreateRemind(h.bot, h.updatesChan, h.curChatID)
	case "1":
		h.SendMsg("1")
	case "2":
		h.SendMsg("2")
	case "3":
		h.SendMsg("3")
	default:
		h.SendMsg("Не знаю таких команд :(")
	}
}
func (h *Handler) SendMsg(message string) {
	h.bot.Send(tgbotapi.NewMessage(h.curUpdate.Message.Chat.ID, message))
}
