package telegram

import (
	"os"
	"os/signal"
	"reminder-bot/internal/database"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	db           database.BotDataBase
	googleApiKey string
	curChatID    int64
	updatesChan  tgbotapi.UpdatesChannel
}

func NewBot(bot *tgbotapi.BotAPI, db database.BotDataBase, googleApiKey string) *Bot {
	return &Bot{
		bot:          bot,
		db:           db,
		googleApiKey: googleApiKey,
	}
}
func (t *Bot) Start() {
	logrus.Infof("Authorized on account %s", t.bot.Self.UserName)

	go t.checkRemindes()

	updates, err := t.initUpdatesChan()
	if err != nil {
		logrus.Warn("Get updates error: ", err)
	}
	t.updatesChan = updates
	go t.rootHandler()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	logrus.Warn("***********Shutdown signal received***********")

	defer t.stopRemindsServing()
}
func (t *Bot) rootHandler() {
	for update := range t.updatesChan {
		if update.Message != nil {
			t.curChatID = update.Message.Chat.ID
			if update.Message.IsCommand() {
				err := t.handleCommand(update.Message)
				if err != nil {
					logrus.Warn("Command handle error: ", err)
				}
			} else {
				err := t.handleMessage(update.Message)
				if err != nil {
					logrus.Warn("Message handle error: ", err)
				}
			}
		}
	}
}

func (t *Bot) initUpdatesChan() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updates, nil
}
func (t *Bot) stopRemindsServing() {
	reminds, err := t.db.GetAllRemindes()
	if err != nil {
		logrus.Warn(err)
	}
	for _, remind := range reminds {
		if remind.ServingNow {
			err = t.db.UpdateRemind(remind, "serving_now", false)
			if err != nil {
				logrus.Warn(err)
			}
		}
	}
	logrus.Warn("***********All gorutines done, shutting down!***********")
}
