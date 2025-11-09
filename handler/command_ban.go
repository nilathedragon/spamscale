package handler

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/restrictions"
)

func CommandBanHandler(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	if matches, err := regexp.Match(`/(\S*) (@[a-zA-Z0-9_-]*)`, []byte(ctx.Message.Text)); err != nil {
		return err
	} else if !matches {
		helperMessage, _ := ctx.Message.Reply(b, "Invalid command format. Use `/ban @user`", &gotgbot.SendMessageOpts{})

		time.AfterFunc(10*time.Second, func() {
			_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
			_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
		})

		return nil
	}

	var mention *gotgbot.MessageEntity
	for _, entity := range ctx.EffectiveMessage.GetEntities() {
		if entity.Type == "mention" ||
			entity.Type == "text_mention" {
			mention = &entity
			break
		}
	}

	if mention == nil {
		helperMessage, err := ctx.Message.Reply(b, "Please add a mention to the user to ban.", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}

		time.AfterFunc(10*time.Second, func() {
			_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
			_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
			return
		})
		return nil
	}

	var userToBan *gotgbot.User
	if mention.Type == "text_mention" {
		userToBan = mention.User
	} else {
		var err error
		userToBan, err = resolveUserFromMention(mention, b, ctx)
		if err != nil {
			return err
		}
	}

	if err := restrictions.Ban(b, userToBan.Id, ctx); err != nil {
		return err
	}
	helperMessage, err := ctx.Message.Reply(b, fmt.Sprintf("Banned %s.", userToBan.Username), &gotgbot.SendMessageOpts{})
	if err != nil {
		return err
	}

	time.AfterFunc(10*time.Second, func() {
		_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
		_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
	})

	return nil
}
