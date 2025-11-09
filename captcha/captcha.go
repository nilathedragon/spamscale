package captcha

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/db"
	"github.com/nilathedragon/spamscale/db/model"
)

func TriggerCaptcha(b *gotgbot.Bot, captchaChatID int64, chatID int64, userID int64) error {
	captchaType, err := db.Chat.GetCaptchaType(chatID)
	if err != nil {
		return err
	}

	log.Info("Triggering captcha", "captcha_type", captchaType, "captcha_chat_id", captchaChatID, "chat_id", chatID, "user_id", userID)
	switch captchaType {
	case model.CaptchaTypeNone:
		return NoneCaptcha(b, captchaChatID, chatID, userID)
	case model.CaptchaTypeButton:
		return ButtonCaptcha(b, captchaChatID, chatID, userID)
	case model.CaptchaTypeEmoji:
		return EmojiCaptcha(b, captchaChatID, chatID, userID)
	case model.CaptchaTypeManual:
		return nil // Manual captcha means the admins will manually process join requests
	default:
		return errors.New("unknown captcha type")
	}
}

func ApproveUser(b *gotgbot.Bot, captchaState *model.CaptchaState) error {
	if captchaState.CaptchaMessageId != -1 {
		if _, err := b.DeleteMessage(captchaState.CaptchaChatId, captchaState.CaptchaMessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
			log.Error("An error occurred while deleting the captcha message", "error", err)
		}
	}

	// Delete the captcha state for this user in this chat
	if err := db.CaptchaState.Delete(captchaState); err != nil {
		log.Error("An error occurred while deleting the captcha state", "error", err)
		return err
	}

	// Approve the user to join the chat
	if _, err := b.ApproveChatJoinRequest(captchaState.ChatID, captchaState.UserID, &gotgbot.ApproveChatJoinRequestOpts{}); err != nil {
		log.Error("An error occurred while approving the chat join request", "error", err)
		return err
	}

	log.Info("User approved", "chat_id", captchaState.ChatID, "user_id", captchaState.UserID)
	return nil
}

func RejectUser(b *gotgbot.Bot, captchaState *model.CaptchaState) error {
	if captchaState.CaptchaMessageId != -1 {
		if _, err := b.DeleteMessage(captchaState.CaptchaChatId, captchaState.CaptchaMessageId, &gotgbot.DeleteMessageOpts{}); err != nil {
			log.Error("An error occurred while deleting the captcha message", "error", err)
		}
	}

	if err := db.CaptchaState.Delete(captchaState); err != nil {
		log.Error("An error occurred while deleting the captcha state", "error", err)
		return err
	}

	if _, err := b.DeclineChatJoinRequest(captchaState.ChatID, captchaState.UserID, &gotgbot.DeclineChatJoinRequestOpts{}); err != nil {
		log.Error("An error occurred while declining the chat join request", "error", err)
		return err
	}

	log.Info("User rejected", "chat_id", captchaState.ChatID, "user_id", captchaState.UserID)
	return nil
}
