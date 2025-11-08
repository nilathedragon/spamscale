package handler

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/captcha"
	"github.com/nilathedragon/spamscale/db"
	"github.com/nilathedragon/spamscale/db/model"
)

func JoinRequestHandlerFilter(cjr *gotgbot.ChatJoinRequest) bool {
	return true
}

func JoinRequestHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	logger := log.With("chat_id", ctx.ChatJoinRequest.Chat.Id, "user_id", ctx.ChatJoinRequest.From.Id)
	logger.Info("Chat join request received")

	captchaType, err := db.Chat.GetCaptchaType(ctx.ChatJoinRequest.Chat.Id)
	if err != nil {
		return err
	}

	switch captchaType {
	case model.CaptchaTypeNone:
		_, err := b.ApproveChatJoinRequest(ctx.ChatJoinRequest.Chat.Id, ctx.ChatJoinRequest.From.Id, &gotgbot.ApproveChatJoinRequestOpts{})
		if err != nil {
			return err
		}
		logger.Info("Chat join request approved, as no captcha was enabled.")
		return nil
	case model.CaptchaTypeButton:
		err := captcha.ButtonCaptcha(b, ctx.ChatJoinRequest.UserChatId, ctx.ChatJoinRequest.Chat.Id, ctx.ChatJoinRequest.From.Id)
		if err != nil {
			return err
		}
		return nil
	case model.CaptchaTypeEmoji:
		err := captcha.EmojiCaptcha(b, ctx.ChatJoinRequest.UserChatId, ctx.ChatJoinRequest.Chat.Id, ctx.ChatJoinRequest.From.Id)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("unknown captcha type")
	}
}
