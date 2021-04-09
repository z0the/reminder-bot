package text

import (
	"fmt"
	"reminder-bot/internal/models"
)

func RemindMessageText(remind *models.Remind) string {
	return fmt.Sprintf(`Ваше напоминание:

				%s

				%s в %s`, remind.Text, remind.ActivationTime.Format("2006.01.02"), remind.ActivationTime.Format("15:04"))
}
