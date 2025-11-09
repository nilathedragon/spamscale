package db

import (
	"fmt"

	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/config"
	"github.com/nilathedragon/spamscale/db/model"
)

var Chat = &chatImpl{}

type chatImpl struct {
	cachedDb
}

const (
	chatCaptchaTypeKey          = "chat:captcha_type:%d"
	chatFastBlocklistEnabledKey = "chat:fast_blocklist_enabled:%d"
)

func (c *chatImpl) GetOrCreate(chatId int64) (*model.Chat, error) {
	var chat model.Chat
	if err := c.getDB().Where(&model.Chat{ID: chatId}).FirstOrCreate(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (c *chatImpl) GetCaptchaType(chatId int64) (model.CaptchaType, error) {
	cacheKey := fmt.Sprintf(chatCaptchaTypeKey, chatId)
	if exists, err := c.getCache().Contains(cacheKey); err != nil {
		return model.CaptchaTypeNone, err
	} else if exists {
		var captchaType model.CaptchaType
		if err := c.getCache().Get(cacheKey, &captchaType); err != nil {
			return model.CaptchaTypeNone, err
		}
		log.Debug("Captcha type found in cache", "chat_id", chatId, "captcha_type", captchaType)
		return captchaType, nil
	}

	var chat model.Chat
	if err := c.getDB().Where(&model.Chat{ID: chatId}).First(&chat).Error; err != nil {
		return model.CaptchaTypeNone, err
	}
	log.Debug("Fetching captcha type for chat", "chat_id", chatId, "captcha_type", chat.CaptchaType)
	c.getCache().Set(cacheKey, chat.CaptchaType)
	c.getCache().ExpireAfter(cacheKey, config.GetCacheDuration())
	return chat.CaptchaType, nil
}

func (c *chatImpl) IsFastBlocklistEnabled(chatId int64) (bool, error) {
	cacheKey := fmt.Sprintf(chatFastBlocklistEnabledKey, chatId)
	if exists, err := c.getCache().Contains(cacheKey); err != nil {
		return false, err
	} else if exists {
		var fastBlocklistEnabled bool
		if err := c.getCache().Get(cacheKey, &fastBlocklistEnabled); err != nil {
			return false, err
		}
		return fastBlocklistEnabled, nil
	}
	var chat model.Chat
	if err := c.getDB().Where(&model.Chat{ID: chatId}).First(&chat).Error; err != nil {
		return false, err
	}
	c.getCache().Set(cacheKey, chat.BlocklistFastEnabled)
	c.getCache().ExpireAfter(cacheKey, config.GetCacheDuration())
	return chat.BlocklistFastEnabled, nil
}

func (c *chatImpl) SetCaptchaType(chatId int64, captchaType model.CaptchaType) error {
	cacheKey := fmt.Sprintf(chatCaptchaTypeKey, chatId)
	log.Debug("Setting captcha type for chat", "chat_id", chatId, "captcha_type", captchaType)
	c.getCache().Set(cacheKey, captchaType)
	c.getCache().ExpireAfter(cacheKey, config.GetCacheDuration())
	return c.getDB().Where(&model.Chat{ID: chatId}).Updates(&model.Chat{CaptchaType: captchaType}).Error
}

func (c *chatImpl) SetFastBlocklistEnabled(chatId int64, fastBlocklistEnabled bool) error {
	cacheKey := fmt.Sprintf(chatFastBlocklistEnabledKey, chatId)
	log.Debug("Setting fast blocklist enabled for chat", "chat_id", chatId, "fast_blocklist_enabled", fastBlocklistEnabled)
	c.getCache().Set(cacheKey, fastBlocklistEnabled)
	c.getCache().ExpireAfter(cacheKey, config.GetCacheDuration())
	return c.getDB().Where(&model.Chat{ID: chatId}).Updates(&model.Chat{BlocklistFastEnabled: fastBlocklistEnabled}).Error
}
