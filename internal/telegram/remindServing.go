package telegram

import (
	"fmt"
	"reminder-bot/internal/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (t *Bot) checkRemindes() {
	logrus.Info("Start reminds serving...")
	for {
		reminds, err := t.db.GetAllRemindes()
		if err != nil {
			logrus.Warn(err)
		}
		for _, remind := range reminds {
			if time.Until(remind.ActivationTime) < time.Hour &&
				!remind.AlreadyServed && !remind.ServingNow {
				go t.serveRemind(remind)
			}
		}
		time.Sleep(55 * time.Minute)
	}
}
func (t *Bot) serveRemind(remind models.Remind) {
	remind.ServingNow = true
	err := t.db.UpdateRemind(remind, "serving_now", remind.ServingNow)
	if err != nil {
		logrus.Warn(err)
	}
	time.Sleep(time.Until(remind.ActivationTime))
	msg := tgbotapi.NewMessage(remind.ChatID, "")
	msg.ParseMode = "HTML"
	msg.Text = fmt.Sprintf(`
	<b>Напоминание!</b>

	%s
	`, remind.Text)
	t.bot.Send(msg)
	remind.ServingNow = false
	remind.AlreadyServed = true
	err = t.db.UpdateRemind(remind, "serving_now", remind.ServingNow)
	if err != nil {
		logrus.Warn(err)
	}
	err = t.db.UpdateRemind(remind, "already_served", remind.AlreadyServed)
	if err != nil {
		logrus.Warn(err)
	}
}
