package handler

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/captcha"
)

func CommandCaptchaHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if _, err := b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
		return err
	}

	return captcha.TriggerCaptcha(b, ctx.EffectiveChat.Id, ctx.EffectiveChat.Id, ctx.EffectiveUser.Id)
}
