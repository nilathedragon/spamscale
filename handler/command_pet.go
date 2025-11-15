package handler

import (
	"bytes"
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/resources"
	"github.com/nilathedragon/spamscale/util"
)

func CommandPetHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	petEntry, ok := petCommandAllowedMap[ctx.Message.Chat.Id]
	if !ok {
		log.Debug("Pet command received, but no pet command was allowed yet.")
		return nil
	}

	if ok, err := util.RateLimitSingle(commandPetUserCacheKey(ctx.Message.Chat.Id, ctx.EffectiveUser.Id), 1*time.Second); err != nil {
		return err
	} else if !ok {
		return nil
	}

	stickerData, err := resources.Read("cute.png")
	if err != nil {
		return err
	}

	msg, err := b.SendSticker(
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
	petEntry.MessagesToDelete = append(petEntry.MessagesToDelete, ctx.Message.MessageId, msg.MessageId)

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

func commandPetUserCacheKey(chatId int64, userId int64) string {
	return fmt.Sprintf("command_pet:%d:%d", chatId, userId)
}
