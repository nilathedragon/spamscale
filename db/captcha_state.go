package db

import (
	"time"

	"github.com/nilathedragon/spamscale/db/model"
)

var CaptchaState = &captchaStateImpl{}

type captchaStateImpl struct {
	cachedDb
}

func (c *captchaStateImpl) Get(captchaChatId int64, captchaMessageId int64, userId int64) (model.CaptchaState, error) {
	var captchaState model.CaptchaState
	if err := c.getDB().Where(&model.CaptchaState{
		CaptchaChatId:    captchaChatId,
		CaptchaMessageId: captchaMessageId,
		UserID:           userId,
	}).First(&captchaState).Error; err != nil {
		return model.CaptchaState{}, err
	}
	return captchaState, nil
}

func (c *captchaStateImpl) Delete(captchaState *model.CaptchaState) error {
	return c.getDB().Delete(captchaState).Error
}

func (c *captchaStateImpl) GetExpiredCaptchas() ([]model.CaptchaState, error) {
	var captchaStates []model.CaptchaState
	if err := c.getDB().Where("expires_at < ?", time.Now()).Find(&captchaStates).Error; err != nil {
		return nil, err
	}
	return captchaStates, nil
}

func (c *captchaStateImpl) Save(captchaState *model.CaptchaState) error {
	return c.getDB().Save(captchaState).Error
}
