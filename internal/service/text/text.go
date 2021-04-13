package text

import (
	"fmt"
	"reminder-bot/internal/models"
	"time"
)

func RemindMessageText(remind *models.Remind, user *models.User) string {
	outTime := remind.ActivationTime.Add(time.Hour * time.Duration(user.TimeZoneOffset))
	return fmt.Sprintf(`Ваше напоминание №%d:

				%s

				%s в %s`, remind.IDForChat, remind.Text, outTime.UTC().Format("2006.01.02"), outTime.UTC().Format("15:04"))
}
func SettingsMessageText(remind *models.Remind) string {
	return fmt.Sprintf(`
	Разница в часах с москвой
	`)
}
