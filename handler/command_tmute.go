package handler

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/restrictions"
	"github.com/nilathedragon/spamscale/util"
)

func CommandTMuteHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !util.IsModerator(b, ctx.EffectiveChat.Id, ctx.EffectiveUser.Id) {
		helperMessage, err := ctx.Message.Reply(b, "You are not authorized to use this command", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}

		queueMessageDeletion(b, []*gotgbot.Message{ctx.Message, helperMessage})
		return nil
	}

	if matches, err := validateCommand(b, ctx.Message.Text, `/(\S*) (@[a-zA-Z0-9_-]*) (\d*(?:min|h|d|m|y))`, "/tmute @user duration", ctx); err != nil {
		return err
	} else if !matches {
		return nil
	}
	command := strings.Split(ctx.Message.Text, " ")

	userToMute, err := getTargetUserFromEntities(b, ctx.Message.GetEntities(), ctx)
	if err != nil {
		return err
	}

	if err := restrictions.TemporaryMute(b, userToMute.Id, command[2], ctx); err != nil {
		return err
	}
	helperMessage, err := ctx.Message.Reply(b, "User muted for "+command[2], &gotgbot.SendMessageOpts{})
	if err != nil {
		return err
	}

	queueMessageDeletion(b, []*gotgbot.Message{ctx.Message, helperMessage})

	return nil
}
