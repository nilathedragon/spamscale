package handler

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/nilathedragon/spamscale/db"
)

const (
	permissionMessage = `Thank you for using Spam Scale in your chat! In order to use the bot, you need to grant me the following admin permissions:
- Invite users
- Delete messages
- Ban users
`
	permissionOkMessage = `Configuration updated, Spam Scale is active!`
)

func SetupHandlerFilter(u *gotgbot.ChatMemberUpdated) bool {
	return true
}

func SetupHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if _, err := db.Chat.GetOrCreate(ctx.MyChatMember.Chat.Id); err != nil {
		return err
	}

	// Check admins of self in chat
	admins, err := b.GetChatAdministrators(ctx.MyChatMember.Chat.Id, &gotgbot.GetChatAdministratorsOpts{})
	if err != nil {
		return err
	}

	if len(admins) != 0 {
		for _, admin := range admins {
			adminUser := admin.MergeChatMember()
			if adminUser.User.Id == ctx.MyChatMember.NewChatMember.GetUser().Id {
				if adminUser.CanDeleteMessages && adminUser.CanRestrictMembers && adminUser.CanInviteUsers {
					_, err = b.SendMessage(ctx.MyChatMember.Chat.Id, permissionOkMessage, &gotgbot.SendMessageOpts{})
					if err != nil {
						return err
					}
					return nil
				}
				break
			}
		}
	}

	_, err = b.SendMessage(ctx.MyChatMember.Chat.Id, permissionMessage, &gotgbot.SendMessageOpts{})
	return err
}
