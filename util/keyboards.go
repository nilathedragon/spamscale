package util

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func GenerateKeyboard[T any](options []T, callbackPrefix string, columns int) [][]gotgbot.InlineKeyboardButton {
	keyboard := make([][]gotgbot.InlineKeyboardButton, 0)
	for i := 0; i < len(options); i += columns {
		row := make([]gotgbot.InlineKeyboardButton, 0)
		for j := 0; j < columns && i+j < len(options); j++ {
			optionText := fmt.Sprintf("%v", options[i+j])
			row = append(row, gotgbot.InlineKeyboardButton{
				Text:         optionText,
				CallbackData: callbackPrefix + optionText,
			})
		}
		keyboard = append(keyboard, row)
	}
	return keyboard
}
