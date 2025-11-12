package handler

import (
	"bytes"
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/resources"
	"github.com/nilathedragon/spamscale/util"
)

func CommandBoopHandler(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	if ok, err := util.RateLimit(commandBoopChatCacheKey(ctx.Message.Chat.Id), 3, 10*time.Second); err != nil {
		return err
	} else if !ok {
		return nil
	}

	if ok, err := util.RateLimitSingle(commandBoopUserCacheKey(ctx.Message.Chat.Id, ctx.EffectiveUser.Id), 10*time.Second); err != nil {
		return err
	} else if !ok {
		return nil
	}

	f, err := resources.Read("booped.png")
	if err != nil {
		return
	}

	if _, err = b.SendSticker(
		ctx.Message.Chat.Id,
		gotgbot.InputFileByReader("booped.png", bytes.NewReader(f)),
		&gotgbot.SendStickerOpts{
			Emoji: "😳",
		},
	); err != nil {
		return
	}

	return
}

func commandBoopChatCacheKey(chatID int64) string {
	return fmt.Sprintf("command_boop:%d", chatID)
}

func commandBoopUserCacheKey(chatID int64, UserID int64) string {
	return fmt.Sprintf("command_boop:%d:%d", chatID, UserID)
}
