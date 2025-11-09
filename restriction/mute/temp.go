package mute

import (
	"errors"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func TemporaryMute(b *gotgbot.Bot, userId int64, duration string, ctx *ext.Context) error {
	regex, err := regexp.Compile(`(\d*)(min|h|d|m|y)`)
	if err != nil {
		return err
	}

	matches := regex.FindStringSubmatch(duration)

	t, err := strconv.Atoi(matches[1])
	t64 := int64(t)
	if err != nil {
		return err
	}
	multiplier := matches[2]

	var deltaTime time.Duration

	switch multiplier {
	case "min":
		deltaTime = time.Duration(t64 * int64(time.Minute))
	case "h":
		deltaTime = time.Duration(t64 * int64(time.Hour))
	case "d":
		deltaTime = time.Duration(t64 * int64(time.Hour*24))
	case "m":
		deltaTime = time.Duration(t64 * int64(time.Hour*24*30))
	case "y":
		deltaTime = time.Duration(t64 * int64(time.Hour*24*365))
	}

	expiration := time.Now().Add(deltaTime)

	member, err := ctx.Message.Chat.GetMember(b, userId, nil)
	if err != nil {
		return err
	}

	admins, err := ctx.Message.Chat.GetAdministrators(b, nil)
	if err != nil {
		return err
	}

	if slices.Contains(admins, member) {
		helperMessage, err := ctx.Message.Reply(b, "You can't mute an adminstrator", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}
		time.AfterFunc(10*time.Second, func() {
			_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
			_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
		})
		return errors.New("user is adminstrator, skipping")
	}

	if _, err := ctx.Message.Chat.RestrictMember(
		b,
		userId,
		gotgbot.ChatPermissions{},
		&gotgbot.RestrictChatMemberOpts{
			UntilDate: expiration.Unix(),
		},
	); err != nil {
		return err
	}

	if _, err := ctx.Message.Chat.BanMember(b, userId, &gotgbot.BanChatMemberOpts{
		UntilDate:      0,
		RevokeMessages: false,
		RequestOpts:    nil,
	}); err != nil {
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
