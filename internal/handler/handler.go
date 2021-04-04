package handler

import (
	"log"
	"reminder-bot/internal/service"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler struct {
	bot          *botapi.BotAPI
	mainKeyboard *botapi.ReplyKeyboardMarkup
	updatesChan  *botapi.UpdatesChannel
	curUpdate    *botapi.Update
	curChatID    int64
	service      *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(bot *botapi.BotAPI, mainKeyboard *botapi.ReplyKeyboardMarkup, updatesChan *botapi.UpdatesChannel) *Handler {
	handler := new(Handler)
	handler.updatesChan = updatesChan
	handler.bot = bot
	handler.mainKeyboard = mainKeyboard
	return handler
}

func (h *Handler) Handle(update *botapi.Update) {
	h.curUpdate = update
	h.curChatID = update.Message.Chat.ID

	switch update.Message.Text {
	case "/start":
		log.Println("Bot is started...")
		msg := botapi.NewMessage(update.Message.Chat.ID, `
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
	h.bot.Send(botapi.NewMessage(h.curUpdate.Message.Chat.ID, message))
}
