package captcha

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/db"
	"github.com/nilathedragon/spamscale/db/model"
)

func NoneCaptcha(b *gotgbot.Bot, captchaChatId int64, chatID int64, userID int64) error {
	log.Info("No captcha was enabled, approving user to join the chat.")

	captchaState := &model.CaptchaState{
		CaptchaMessageId: -1,
		CaptchaChatId:    captchaChatId,
		ChatID:           chatID,
		UserID:           userID,
		ExpiresAt:        time.Now().Add(time.Minute * 10),
	}

	if err := db.CaptchaState.Save(captchaState); err != nil {
		return err
	}

	return ApproveUser(b, captchaState)
}
