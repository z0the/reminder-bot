package keyboard

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetMainKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Создать напоминание"),
	}
	row2 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Посмотреть список моих напоминаний"),
		tgbotapi.NewKeyboardButton("Настройки"),
	}
	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)
	return &keyboard
}

// Two buttons "Отмена" and "Далее"
func GetCancelNextKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("Отмена"), tgbotapi.NewKeyboardButton("Далее")})
	return &keyboard
}

// Two buttons "-" and "+"
func GetSettingTimeOffsetKeyboard() *tgbotapi.InlineKeyboardMarkup {
	row1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("+1", "+"),
		tgbotapi.NewInlineKeyboardButtonData("-1", "-"),
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1)
	return &keyboard
}
func GetDeleteKeyboard(id uint) *tgbotapi.InlineKeyboardMarkup {
	row1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Удалить", fmt.Sprint(id)),
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1)
	return &keyboard
}

func GetBackwardKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("Назад")})
	return &keyboard
}
func GetBackwardSaveKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Назад"),
		tgbotapi.NewKeyboardButton("Сохранить"),
	}
	keyboard := tgbotapi.NewReplyKeyboard(row1)
	return &keyboard
}
