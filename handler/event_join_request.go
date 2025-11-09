package handler

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/captcha"
	"github.com/nilathedragon/spamscale/db"
	"github.com/nilathedragon/spamscale/pkg/thirdparty/fast"
)

func JoinRequestHandlerFilter(cjr *gotgbot.ChatJoinRequest) bool {
	return true
}

func JoinRequestHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	logger := log.With("chat_id", ctx.ChatJoinRequest.Chat.Id, "user_id", ctx.ChatJoinRequest.From.Id)
	logger.Info("Chat join request received")

	if fastBlocklistEnabled, err := db.Chat.IsFastBlocklistEnabled(ctx.ChatJoinRequest.Chat.Id); err != nil {
		logger.Error("Failed to check if fast blocklist is enabled", "error", err)
	} else if fastBlocklistEnabled {
		if blocked, err := fast.IsBlocked(ctx.ChatJoinRequest.From.Id); err != nil {
			return err
		} else if blocked {
			logger.Info("User is blocked by fast blocklist, rejecting join request")
			if _, err := b.DeclineChatJoinRequest(ctx.ChatJoinRequest.Chat.Id, ctx.ChatJoinRequest.From.Id, &gotgbot.DeclineChatJoinRequestOpts{}); err != nil {
				return err
			}
			return nil
		}
	}

	return captcha.TriggerCaptcha(b, ctx.ChatJoinRequest.UserChatId, ctx.ChatJoinRequest.Chat.Id, ctx.ChatJoinRequest.From.Id)
}
