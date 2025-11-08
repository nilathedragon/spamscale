package util

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func IsAdmin(b *gotgbot.Bot, chatID int64, userID int64) bool {
	admins, err := b.GetChatAdministrators(chatID, &gotgbot.GetChatAdministratorsOpts{})
	if err != nil {
		return false
	}
	for _, admin := range admins {
		adminUser := admin.MergeChatMember()
		if adminUser.User.Id == userID {
			return adminUser.CanPromoteMembers || admin.GetStatus() == "creator"
		}
	}

	return false
}

func IsModerator(b *gotgbot.Bot, chatID int64, userID int64) bool {
	admins, err := b.GetChatAdministrators(chatID, &gotgbot.GetChatAdministratorsOpts{})
	if err != nil {
		return false
	}
	for _, admin := range admins {
		adminUser := admin.MergeChatMember()
		if adminUser.User.Id == userID {
			return adminUser.CanRestrictMembers || admin.GetStatus() == "creator"
		}
	}
	return false
}
