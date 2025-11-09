package handler

import (
	"regexp"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/helpers"
	"github.com/nilathedragon/spamscale/restriction/mute"
)

type TgUserData struct {
	ID         int64  `json:"id,string"`
	Type       string `json:"type"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Username   string `json:"username"`
	IsBot      bool   `json:"isBot"`
	IsVerified bool   `json:"isVerified"`
}

func CommandTMuteHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if matches, err := regexp.Match(`/(\S*) (@[a-zA-Z0-9_-]*) (\d*(?:min|h|d|m|y))`, []byte(ctx.Message.Text)); err != nil {
		return err
	} else if !matches {
		helperMessage, _ := ctx.Message.Reply(b, "Invalid command format. Use `/tmute @user duration`", &gotgbot.SendMessageOpts{})

		time.AfterFunc(10*time.Second, func() {
			_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
			_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
		})

		return nil
	}
	command := strings.Split(ctx.Message.Text, " ")

	var mention *gotgbot.MessageEntity
	for _, entity := range ctx.EffectiveMessage.GetEntities() {
		if entity.Type == "mention" ||
			entity.Type == "text_mention" {
			mention = &entity
			break
		}
	}

	if mention == nil {
		helperMessage, err := ctx.Message.Reply(b, "Please add a mention to the user to mute.", &gotgbot.SendMessageOpts{})
		if err != nil {
			return err
		}

		time.AfterFunc(10*time.Second, func() {
			_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
			_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
			return
		})
		return nil
	}

	var userToMute *gotgbot.User
	if mention.Type == "text_mention" {
		userToMute = mention.User
	} else {
		var err error
		userToMute, err = resolveUserFromMention(mention, b, ctx)
		if err != nil {
			return err
		}
	}

	if err := mute.TemporaryMute(b, ctx.EffectiveChat.Id, userToMute.Id, command[2], ctx); err != nil {
		return err
	}
	helperMessage, err := ctx.Message.Reply(b, "User muted for "+command[2], &gotgbot.SendMessageOpts{})
	if err != nil {
		return err
	}

	time.AfterFunc(10*time.Second, func() {
		_, _ = ctx.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
		_, _ = helperMessage.Delete(b, &gotgbot.DeleteMessageOpts{})
	})

	return nil
}

func resolveUserFromMention(mention *gotgbot.MessageEntity, b *gotgbot.Bot, ctx *ext.Context) (*gotgbot.User, error) {
	username := ctx.Message.Text[mention.Offset+1 : mention.Offset+mention.Length]

	userID, err := helpers.GetUserIDFromUsername(username)
	if err != nil {
		return nil, err
	}

	member, err := ctx.Message.Chat.GetMember(b, userID, &gotgbot.GetChatMemberOpts{})
	if err != nil {
		return nil, err
	}

	usr := member.GetUser()

	return &usr, nil
}
