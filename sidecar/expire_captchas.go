package sidecar

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/captcha"
	"github.com/nilathedragon/spamscale/db"
)

func ExpireCaptchas(b *gotgbot.Bot) {
	for {
		expireCaptchas(b)
		<-time.After(time.Minute)
	}
}

func expireCaptchas(b *gotgbot.Bot) {
	expiredCaptchas, err := db.CaptchaState.GetExpiredCaptchas()
	if err != nil {
		log.Error("Failed to get expired captchas", "error", err)
		return
	}

	for _, captchaState := range expiredCaptchas {
		captcha.RejectUser(b, &captchaState)
		db.CaptchaState.Delete(&captchaState)
	}
}
