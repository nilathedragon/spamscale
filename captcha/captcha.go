package captcha

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/go-mojito/mojito/log"
	"github.com/infinytum/injector"
	"github.com/nilathedragon/spamscale/db/model"
	"gorm.io/gorm"
)

func ApproveUser(b *gotgbot.Bot, captchaState *model.CaptchaState) error {
	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		log.Error("An error occurred while injecting the database", "error", err)
		return err
	}

	// Delete the captcha state for this user in this chat
	if err := db.Delete(captchaState).Error; err != nil {
		log.Error("An error occurred while deleting the captcha state", "error", err)
		return err
	}

	// Approve the user to join the chat
	_, err = b.ApproveChatJoinRequest(captchaState.ChatID, captchaState.UserID, &gotgbot.ApproveChatJoinRequestOpts{})
	if err != nil {
		return err
	} else {
		log.Info("User approved", "chat_id", captchaState.ChatID, "user_id", captchaState.UserID)
	}

	return nil
}

func RejectUser(b *gotgbot.Bot, captchaState *model.CaptchaState) error {
	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		log.Error("An error occurred while injecting the database", "error", err)
		return err
	}

	if err := db.Delete(captchaState).Error; err != nil {
		log.Error("An error occurred while deleting the captcha state", "error", err)
		return err
	}

	_, err = b.DeclineChatJoinRequest(captchaState.ChatID, captchaState.UserID, &gotgbot.DeclineChatJoinRequestOpts{})
	if err != nil {
		return err
	} else {
		log.Info("User rejected", "chat_id", captchaState.ChatID, "user_id", captchaState.UserID)
	}

	return nil
}
