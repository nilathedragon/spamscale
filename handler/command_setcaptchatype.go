package handler

import (
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/infinytum/injector"
	"github.com/nilathedragon/spamscale/db/model"
	"github.com/nilathedragon/spamscale/util"
	"gorm.io/gorm"
)

const (
	CommandSetCaptchaTypeCallback = "set_captcha_type:"
)

func CommandSetCaptchaTypeHandler(b *gotgbot.Bot, ctx *ext.Context) error {
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

	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		return err
	}

	var chat model.Chat
	if err := db.Where(&model.Chat{ID: ctx.EffectiveChat.Id}).First(&chat).Error; err != nil {
		return err
	}
	chat.CaptchaType = model.CaptchaType(captchaType)
	if err := db.Save(&chat).Error; err != nil {
		return err
	}

	_, err = b.SendMessage(ctx.EffectiveChat.Id, "Captcha type updated successfully", &gotgbot.SendMessageOpts{})
	if err != nil {
		return err
	}

	return nil
}
