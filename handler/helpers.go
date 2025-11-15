package handler

import (
	"errors"
	"regexp"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/pkg/thirdparty/tguserid"
)

func validateCommand(b *gotgbot.Bot, cmd string, regex string, expectedFormat string, ctx *ext.Context) (bool, error) {
	if matches, err := regexp.Match(regex, []byte(cmd)); err != nil {
		return false, err
	} else if !matches {
		helperMessage, _ := ctx.Message.Reply(b, "Invalid command format. Use `"+expectedFormat+"`", &gotgbot.SendMessageOpts{})
		queueMessageDeletion(b, []*gotgbot.Message{ctx.Message, helperMessage})
		return false, nil
	}

	return true, nil
}

func queueMessageDeletion(b *gotgbot.Bot, msgs []*gotgbot.Message) {
	time.AfterFunc(10*time.Second, func() {
		for _, msg := range msgs {
			_, _ = msg.Delete(b, &gotgbot.DeleteMessageOpts{})
		}
	})
}

func getTargetUserFromEntities(b *gotgbot.Bot, entities []gotgbot.MessageEntity, ctx *ext.Context) (*gotgbot.User, error) {
	var mention *gotgbot.MessageEntity
	for _, entity := range entities {
		if entity.Type == "mention" ||
			entity.Type == "text_mention" {
			mention = &entity
			break
		}
	}

	if mention == nil {
		helperMessage, err := ctx.Message.Reply(b, "Please add a mention to the user to target.", &gotgbot.SendMessageOpts{})
		if err != nil {
			return nil, err
		}

		queueMessageDeletion(b, []*gotgbot.Message{helperMessage, ctx.Message})
		return nil, errors.New("no mention found")
	}

	var user *gotgbot.User
	if mention.Type == "text_mention" {
		user = mention.User
	} else {
		var err error
		user, err = resolveUserFromMention(mention, b, ctx)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func resolveUserFromMention(mention *gotgbot.MessageEntity, b *gotgbot.Bot, ctx *ext.Context) (*gotgbot.User, error) {
	username := ctx.Message.Text[mention.Offset+1 : mention.Offset+mention.Length]

	user, err := tguserid.GetUser(username)
	if err != nil {
		return nil, err
	}

	member, err := ctx.Message.Chat.GetMember(b, user.ID, &gotgbot.GetChatMemberOpts{})
	if err != nil {
		return nil, err
	}

	usr := member.GetUser()

	return &usr, nil
}
