package restrictions

import (
	"errors"
	"slices"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Ban(b *gotgbot.Bot, userId int64, ctx *ext.Context) error {
	member, err := ctx.Message.Chat.GetMember(b, userId, nil)
	if err != nil {
		return err
	}

	admins, err := ctx.Message.Chat.GetAdministrators(b, nil)
	if err != nil {
		return err
	}

	if slices.Contains(admins, member) {
		helperMessage, err := ctx.Message.Reply(b, "You can't ban an adminstrator", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}
		time.AfterFunc(10*time.Second, func() {
			_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
			_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
		})
		return errors.New("user is adminstrator, skipping")
	}

	if _, err := ctx.Message.Chat.BanMember(b, userId, nil); err != nil {
		helperMessage, err := ctx.Message.Reply(b, "Failed to ban user", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}

		time.AfterFunc(10*time.Second, func() {
			_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
			_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
		})

		return err
	}

	return nil
}
