package util

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func TempMessage(b *gotgbot.Bot, chatId int64, message string) error {
	msg, err := b.SendMessage(chatId, message, &gotgbot.SendMessageOpts{})
	if err != nil {
		return err
	}
	time.AfterFunc(10*time.Second, func() {
		_, _ = msg.Delete(b, &gotgbot.DeleteMessageOpts{})
	})
	return nil
}

func DropMessage(msg *gotgbot.Message, err error) error {
	return err
}
