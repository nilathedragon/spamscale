package util

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func IsAdmin(b *gotgbot.Bot, chatID int64, userID int64) bool {
	member, err := b.GetChatMember(chatID, userID, &gotgbot.GetChatMemberOpts{})
	if err != nil {
		return false
	}
	return isAdmin(member)
}

func IsModerator(b *gotgbot.Bot, chatID int64, userID int64) bool {
	member, err := b.GetChatMember(chatID, userID, &gotgbot.GetChatMemberOpts{})
	if err != nil {
		return false
	}
	return isModerator(member)
}

func GetModerators(b *gotgbot.Bot, chatID int64) ([]gotgbot.ChatMember, error) {
	admins, err := b.GetChatAdministrators(chatID, &gotgbot.GetChatAdministratorsOpts{})
	if err != nil {
		return nil, err
	}

	moderators := make([]gotgbot.ChatMember, 0)
	for _, admin := range admins {
		if admin.GetUser().IsBot {
			continue
		}

		if isModerator(admin.MergeChatMember()) {
			moderators = append(moderators, admin.MergeChatMember())
		}
	}
	return moderators, nil
}

func isModerator(member gotgbot.ChatMember) bool {
	return member.MergeChatMember().CanRestrictMembers || member.GetStatus() == "creator"
}

func isAdmin(member gotgbot.ChatMember) bool {
	return member.MergeChatMember().CanPromoteMembers || member.GetStatus() == "creator"
}
