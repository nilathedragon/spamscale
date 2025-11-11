package handler

import (
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func CommandPetHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	f, err := os.Open("./res/cute.png")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = b.SendSticker(
		ctx.Message.Chat.Id,
		&gotgbot.FileReader{
			Name: "cute.png",
			Data: f,
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
