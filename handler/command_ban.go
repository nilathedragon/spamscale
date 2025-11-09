package handler

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/restrictions"
)

func CommandBanHandler(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	if valid, err := validateCommand(b, ctx.Message.Text, `/(\S*) (@[a-zA-Z0-9_-]*)`, "/ban @user", ctx); err != nil {
		return err
	} else if !valid {
		return nil
	}

	userToBan, err := getTargetUserFromEntities(b, ctx.EffectiveMessage.GetEntities(), ctx)
	if err != nil {
		return err
	}

	if err := restrictions.Ban(b, userToBan.Id, ctx); err != nil {
		return err
	}
	helperMessage, err := ctx.Message.Reply(b, fmt.Sprintf("Banned %s.", userToBan.Username), &gotgbot.SendMessageOpts{})
	if err != nil {
		return err
	}

	queueMessageDeletion(b, []*gotgbot.Message{ctx.Message, helperMessage})

	return nil
}
