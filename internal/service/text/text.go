package text

import (
	"fmt"
	"reminder-bot/internal/models"
)

func RemindMessageText(remind *models.Remind) string {
	return fmt.Sprintf(`Ваше напоминание №%d:

				%s

				%s в %s`, remind.IDForChat, remind.Text, remind.ActivationTime.Format("2006.01.02"), remind.ActivationTime.Format("15:04"))
}
func SettingsMessageText(remind *models.Remind) string {
	return fmt.Sprintf(`
	Разница в часах с москвой
	`,)
}