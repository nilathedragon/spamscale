package handler

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/restrictions"
	"github.com/nilathedragon/spamscale/util"
)

func CommandTMuteHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Info("TMute command received", "chat_id", ctx.EffectiveChat.Id, "user_id", ctx.EffectiveUser.Id)
	if !util.IsModerator(b, ctx.EffectiveChat.Id, ctx.EffectiveUser.Id) {
		return util.TempMessage(b, ctx.EffectiveChat.Id, "You are not authorized to use this command")
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

	return util.TempMessage(b, ctx.EffectiveChat.Id, "User muted for "+command[2])
}
