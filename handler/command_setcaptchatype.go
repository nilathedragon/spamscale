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
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "You are not authorized to use this command", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}
		return nil
	}

	keyboard := util.GenerateKeyboard(model.CaptchaTypes, CommandSetCaptchaTypeCallback, 2)
	_, err := b.SendMessage(ctx.EffectiveChat.Id, "Please select the type of captcha you want to use", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func CommandSetCaptchaTypeHandlerCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	if !util.IsAdmin(b, ctx.EffectiveChat.Id, ctx.EffectiveUser.Id) {
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "You are not authorized to use this command", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}
		return nil
	}

	if _, err := b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
		return err
	}

	captchaType := ctx.CallbackQuery.Data[len(CommandSetCaptchaTypeCallback):]
	if !slices.Contains(model.CaptchaTypes, model.CaptchaType(captchaType)) {
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "Invalid captcha type. Please try again.", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}
		return nil
	}

	if err := db.Chat.SetCaptchaType(ctx.EffectiveChat.Id, model.CaptchaType(captchaType)); err != nil {
		return err
	}

	if _, err := b.SendMessage(ctx.EffectiveChat.Id, "Captcha type updated successfully", &gotgbot.SendMessageOpts{}); err != nil {
		return err
	}
	return nil
}
