package telegram

import (
	"reminder-bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	db  database.BotDataBase
}

func NewBot(bot *tgbotapi.BotAPI, db database.BotDataBase) *Bot {
	return &Bot{
		bot: bot,
		db:  db,
	}
}
func (b *Bot) Start() {
	logrus.Infof("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChan()
	if err != nil {
		logrus.Warn("Get updates error: ", err)
	}
	b.handleUpdatesChan(updates)
}
func (b *Bot) handleUpdatesChan(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				err := b.handleCommand(update.Message)
				if err != nil {
					logrus.Warn("Command handle error: ", err)
				}
			} else {
				b.handleMessage(update.Message)
			}
		}
	}
}

func (b *Bot) initUpdatesChan() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updates, nil
}
