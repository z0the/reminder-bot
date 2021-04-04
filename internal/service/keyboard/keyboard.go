package keyboard

import botapi "github.com/go-telegram-bot-api/telegram-bot-api"

var (
	mainKeyBoard *botapi.ReplyKeyboardMarkup = nil
)

func GetMainKeyboard() *botapi.ReplyKeyboardMarkup {
	if mainKeyBoard == nil {
		row1 := []botapi.KeyboardButton{
			botapi.NewKeyboardButton("Создать напоминание"),
		}
		row2 := []botapi.KeyboardButton{
			botapi.NewKeyboardButton("Посмотреть список моих напоминаний"),
		}
		row3 := []botapi.KeyboardButton{
			botapi.NewKeyboardButton("1"),
			botapi.NewKeyboardButton("2"),
			botapi.NewKeyboardButton("3"),
		}
		keyboard := botapi.NewReplyKeyboard(row1, row2, row3)
		mainKeyBoard = &keyboard
	}
	return mainKeyBoard
}
