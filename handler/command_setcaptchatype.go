package handler

import (
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/db"
	"github.com/nilathedragon/spamscale/db/model"
	"github.com/nilathedragon/spamscale/util"
)

const (
	CommandSetCaptchaTypeCallback = "set_captcha_type:"
)

func CommandSetCaptchaTypeHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Info("Set captcha type command received", "chat_id", ctx.EffectiveChat.Id, "user_id", ctx.EffectiveUser.Id)
	if _, err := b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
		return err
	}

	if !util.IsAdmin(b, ctx.EffectiveChat.Id, ctx.EffectiveUser.Id) {
		return util.TempMessage(b, ctx.EffectiveChat.Id, "You are not authorized to use this command")
	}

	return util.DropMessage(b.SendMessage(ctx.EffectiveChat.Id, "Please select the type of captcha you want to use", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: util.GenerateKeyboard(model.CaptchaTypes, CommandSetCaptchaTypeCallback, 2),
		},
	}))
}

func CommandSetCaptchaTypeHandlerCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	if !util.IsAdmin(b, ctx.EffectiveChat.Id, ctx.EffectiveUser.Id) {
		return util.TempMessage(b, ctx.EffectiveChat.Id, "You are not authorized to use this command")
	}

	if _, err := b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
		return err
	}

	captchaType := ctx.CallbackQuery.Data[len(CommandSetCaptchaTypeCallback):]
	if !slices.Contains(model.CaptchaTypes, model.CaptchaType(captchaType)) {
		return util.TempMessage(b, ctx.EffectiveChat.Id, "Invalid captcha type. Please try again.")
	}

	if err := db.Chat.SetCaptchaType(ctx.EffectiveChat.Id, model.CaptchaType(captchaType)); err != nil {
		return err
	}

	return util.TempMessage(b, ctx.EffectiveChat.Id, "Captcha type updated successfully")
}
