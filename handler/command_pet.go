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

func CommandPetHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if ok, err := util.RateLimit(commandPetChatCacheKey(ctx.Message.Chat.Id), 3, 10*time.Second); err != nil {
		return err
	} else if !ok {
		return nil
	}

	if ok, err := util.RateLimitSingle(commandPetUserCacheKey(ctx.Message.Chat.Id, ctx.EffectiveUser.Id), 10*time.Second); err != nil {
		return err
	} else if !ok {
		return nil
	}

	stickerData, err := resources.Read("cute.png")
	if err != nil {
		return err
	}

	_, err = b.SendSticker(
		ctx.Message.Chat.Id,
		&gotgbot.FileReader{
			Name: "cute.png",
			Data: bytes.NewReader(stickerData),
		},
		&gotgbot.SendStickerOpts{
			Emoji: "🩵",
		})
	if err != nil {
		return err
	}

	_, err = ctx.Message.SetReaction(b, &gotgbot.SetMessageReactionOpts{
		Reaction: []gotgbot.ReactionType{
			gotgbot.ReactionTypeEmoji{
				Emoji: "❤️",
			},
		},
		IsBig:       true,
		RequestOpts: nil,
	})

	return err
}

func commandPetChatCacheKey(chatId int64) string {
	return fmt.Sprintf("command_pet:%d", chatId)
}

func commandPetUserCacheKey(chatId int64, userId int64) string {
	return fmt.Sprintf("command_pet:%d:%d", chatId, userId)
}
