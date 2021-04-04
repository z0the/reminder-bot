package clock

import (
	"strconv"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GenerateClockKeyboard(hour10, hour1, minute10, minute1 int) botapi.InlineKeyboardMarkup {
	row1 := []botapi.InlineKeyboardButton{
		botapi.NewInlineKeyboardButtonData("+", "+hour10"),
		botapi.NewInlineKeyboardButtonData("+", "+hour1"),
		botapi.NewInlineKeyboardButtonData(" ", "0"),
		botapi.NewInlineKeyboardButtonData("+", "+minute10"),
		botapi.NewInlineKeyboardButtonData("+", "+minute1"),
	}
	row2 := []botapi.InlineKeyboardButton{
		botapi.NewInlineKeyboardButtonData(strconv.Itoa(hour10), "0"),
		botapi.NewInlineKeyboardButtonData(strconv.Itoa(hour1), "0"),
		botapi.NewInlineKeyboardButtonData(":", "0"),
		botapi.NewInlineKeyboardButtonData(strconv.Itoa(minute10), "0"),
		botapi.NewInlineKeyboardButtonData(strconv.Itoa(minute1), "0"),
	}
	row3 := []botapi.InlineKeyboardButton{
		botapi.NewInlineKeyboardButtonData("-", "-hour10"),
		botapi.NewInlineKeyboardButtonData("-", "-hour1"),
		botapi.NewInlineKeyboardButtonData(" ", "0"),
		botapi.NewInlineKeyboardButtonData("-", "-minute10"),
		botapi.NewInlineKeyboardButtonData("-", "-minute1"),
	}
	return botapi.NewInlineKeyboardMarkup(row1, row2, row3)
}
func PlusHandler(data string, clock map[string]int) (botapi.InlineKeyboardMarkup, map[string]int) {
	switch data {
	case "+hour10":
		if clock["hour10"] < 2 {
			clock["hour10"] += 1
			if clock["hour10"] == 2 && clock["hour1"] > 3 {
				clock["hour1"] = 0
			}
		} else {
			clock["hour10"] = 0
		}
	case "+hour1":
		if clock["hour10"] == 2 {
			if clock["hour1"] < 3 {
				clock["hour1"] += 1
			} else {
				clock["hour1"] = 0
			}
		} else if clock["hour10"] < 2 {
			if clock["hour1"] < 9 {
				clock["hour1"] += 1
			} else {
				clock["hour1"] = 0
			}
		}
	case "+minute10":
		if clock["minute10"] < 5 {
			clock["minute10"] += 1
		} else {
			clock["minute10"] = 0
		}
	case "+minute1":
		if clock["minute1"] < 9 {
			clock["minute1"] += 1
		} else {
			clock["minute1"] = 0
		}
	}
	return GenerateClockKeyboard(clock["hour10"], clock["hour1"], clock["minute10"], clock["minute1"]), clock
}
func MinusHandler(data string, clock map[string]int) (botapi.InlineKeyboardMarkup, map[string]int) {
	switch data {
	case "-hour10":
		if clock["hour10"] == 0 {
			clock["hour10"] = 2
			if clock["hour1"] > 3 {
				clock["hour1"] = 0
			}
		} else {
			clock["hour10"] -= 1
		}
	case "-hour1":
		if clock["hour10"] == 2 {
			if clock["hour1"] == 0 {
				clock["hour1"] = 3
			} else {
				clock["hour1"] -= 1
			}
		} else if clock["hour10"] < 2 {
			if clock["hour1"] == 0 {
				clock["hour1"] = 9
			} else {
				clock["hour1"] -= 1
			}
		}
	case "-minute10":
		if clock["minute10"] == 0 {
			clock["minute10"] = 5
		} else {
			clock["minute10"] -= 1
		}
	case "-minute1":
		if clock["minute1"] == 0 {
			clock["minute1"] = 9
		} else {
			clock["minute1"] -= 1
		}
	}
	return GenerateClockKeyboard(clock["hour10"], clock["hour1"], clock["minute10"], clock["minute1"]), clock
}
