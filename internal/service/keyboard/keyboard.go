package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var (
	mainKeyBoard *tgbotapi.ReplyKeyboardMarkup = nil
)

func GetMainKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	if mainKeyBoard == nil {
		row1 := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Создать напоминание"),
		}
		row2 := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Посмотреть список моих напоминаний"),
		}
		row3 := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("1"),
			tgbotapi.NewKeyboardButton("2"),
			tgbotapi.NewKeyboardButton("3"),
		}
		keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
		mainKeyBoard = &keyboard
	}
	return mainKeyBoard
}
func GetCancelNextKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("Отмена"), tgbotapi.NewKeyboardButton("Далее")})
	return &keyboard
}
func GetSettingTimeOffsetKeyboard() *tgbotapi.InlineKeyboardMarkup {
	row1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("+1", "+"),
		tgbotapi.NewInlineKeyboardButtonData("-1", "-"),
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1)
	return &keyboard
}
