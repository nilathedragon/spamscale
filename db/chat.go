package db

import "github.com/nilathedragon/spamscale/db/model"

var Chat = &chatImpl{}

type chatImpl struct {
	cachedDb
}

func (c *chatImpl) GetCaptchaType(chatId int64) (model.CaptchaType, error) {
	var chat model.Chat
	if err := c.getDB().Where(&model.Chat{ID: chatId}).First(&chat).Error; err != nil {
		return model.CaptchaTypeNone, err
	}
	return chat.CaptchaType, nil
}
