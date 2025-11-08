package handler

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/captcha"
)

func JoinRequestHandlerFilter(cjr *gotgbot.ChatJoinRequest) bool {
	return true
}

func JoinRequestHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	logger := log.With("chat_id", ctx.ChatJoinRequest.Chat.Id, "user_id", ctx.ChatJoinRequest.From.Id)
	logger.Info("Chat join request received")
	return captcha.TriggerCaptcha(b, ctx.ChatJoinRequest.UserChatId, ctx.ChatJoinRequest.Chat.Id, ctx.ChatJoinRequest.From.Id)
}
