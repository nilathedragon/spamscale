package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/util"
)

func CommandReportHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Info("Report command received", "chat_id", ctx.EffectiveChat.Id, "user_id", ctx.EffectiveUser.Id)
	// We always want to delete the message, no matter if it triggers a mod ping or not
	if _, err := b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
		return err
	}

	// Debounce logic, to prevent spamming the mods (unintentionally or intentionally)
	if ok, err := util.RateLimitSingle(commandReportCacheKey(ctx.EffectiveChat.Id), 10*time.Second); err != nil {
		return err
	} else if !ok {
		return nil
	}

	moderators, err := util.GetModerators(b, ctx.EffectiveChat.Id)
	if err != nil {
		return err
	}

	// This little unicode hack allows us to send a message that seemingly mentions no one, but still causes a mention notification to be sent to the mods
	messageEntities := make([]string, len(moderators))
	for i, moderator := range moderators {
		moderatorUser := moderator.GetUser()
		messageEntities[i] = fmt.Sprintf("[\u200b](tg://user?id=%d)", moderatorUser.Id)
	}

	_, err = b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("Thank you for reporting, the moderators have been notified. %s", strings.Join(messageEntities, " ")), &gotgbot.SendMessageOpts{
		ParseMode: "Markdown",
	})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func commandReportCacheKey(chatId int64) string {
	return fmt.Sprintf("command_report:%d", chatId)
}
