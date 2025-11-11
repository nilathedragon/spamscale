package handler

import (
	"bytes"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/resources"
)

func CommandPetHandler(b *gotgbot.Bot, ctx *ext.Context) error {
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
