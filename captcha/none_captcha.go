package captcha

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/go-mojito/mojito/log"
	"github.com/infinytum/injector"
	"github.com/nilathedragon/spamscale/db/model"
	"gorm.io/gorm"
)

func NoneCaptcha(b *gotgbot.Bot, captchaChatId int64, chatID int64, userID int64) error {
	log.Info("No captcha was enabled, approving user to join the chat.")

	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		return err
	}

	captchaState := &model.CaptchaState{
		CaptchaMessageId: 0,
		CaptchaChatId:    captchaChatId,
		ChatID:           chatID,
		UserID:           userID,
		ExpiresAt:        time.Now().Add(time.Minute * 10),
	}

	if err := db.Create(captchaState).Error; err != nil {
		return err
	}

	return ApproveUser(b, captchaState)
}
