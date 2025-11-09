package handler

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/db"
	"github.com/nilathedragon/spamscale/util"
)

const (
	CommandSetFastBlocklistEnabledCallback = "set_fast_blocklist_enabled:"
)

func CommandSetFastBlocklistEnabledHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Info("Set fast blocklist enabled command received", "chat_id", ctx.EffectiveChat.Id, "user_id", ctx.EffectiveUser.Id)
	if _, err := b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
		return err
	}

	if !util.IsAdmin(b, ctx.EffectiveChat.Id, ctx.EffectiveUser.Id) {
		return util.TempMessage(b, ctx.EffectiveChat.Id, "You are not authorized to use this command")
	}

	return util.DropMessage(b.SendMessage(ctx.EffectiveChat.Id, "Please select whether to enable the fast blocklist", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: util.GenerateKeyboard([]string{"Yes", "No"}, CommandSetFastBlocklistEnabledCallback, 2),
		},
	}))
}

func CommandSetFastBlocklistEnabledHandlerCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	if !util.IsAdmin(b, ctx.EffectiveChat.Id, ctx.EffectiveUser.Id) {
		return util.TempMessage(b, ctx.EffectiveChat.Id, "You are not authorized to use this command")
	}

	if _, err := b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
		return err
	}

	fastBlocklistEnabled := ctx.CallbackQuery.Data[len(CommandSetFastBlocklistEnabledCallback):]
	if err := db.Chat.SetFastBlocklistEnabled(ctx.EffectiveChat.Id, fastBlocklistEnabled == "Yes"); err != nil {
		return err
	}

	return util.TempMessage(b, ctx.EffectiveChat.Id, "Fast blocklist setting updated successfully")
}
