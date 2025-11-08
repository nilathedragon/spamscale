package captcha

import (
	"errors"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/infinytum/injector"
	"github.com/nilathedragon/spamscale/db/model"
	"gorm.io/gorm"
)

const ButtonCaptchaConfirmCallback = "button_captcha_confirm"

func ButtonCaptcha(b *gotgbot.Bot, captchaChatId int64, chatID int64, userID int64) error {
	message, err := b.SendMessage(captchaChatId, "Please confirm you are not a robot by clicking the button below", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Confirm",
						CallbackData: ButtonCaptchaConfirmCallback,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		return err
	}

	return db.Create(&model.CaptchaState{
		CaptchaMessageId: message.MessageId,
		CaptchaChatId:    captchaChatId,
		ChatID:           chatID,
		UserID:           userID,
		ExpiresAt:        time.Now().Add(time.Minute * 10),
	}).Error
}

func ButtonCaptchaCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	messageId := ctx.CallbackQuery.Message.GetMessageId()
	userId := ctx.CallbackQuery.From.Id
	chatId := ctx.CallbackQuery.Message.GetChat().Id

	log.Info("Button captcha callback received", "message_id", messageId, "user_id", userId, "chat_id", chatId)

	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		log.Error("An error occurred while injecting the database", "error", err)
		ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "An error occurred while processing the captcha",
			ShowAlert: true,
		})
		return err
	}

	var captchaState model.CaptchaState
	if err := db.Where(&model.CaptchaState{
		CaptchaMessageId: messageId,
		CaptchaChatId:    chatId,
		UserID:           userId,
	}).First(&captchaState).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Captcha state not found, deleting message", "message_id", messageId, "user_id", userId, "chat_id", chatId)
			_, err = ctx.CallbackQuery.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
			return err
		}

		log.Error("An error occurred while finding the captcha state", "error", err)
		ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "An error occurred while processing the captcha",
			ShowAlert: true,
		})
		return err
	}

	_, err = ctx.CallbackQuery.Message.Delete(b, &gotgbot.DeleteMessageOpts{})
	if err != nil {
		log.Error("An error occurred while deleting the captcha message", "error", err)
	}

	// Retry the captcha if it has expired, but the user somehow managed to catch it before the cleanup cronjob ran
	if captchaState.ExpiresAt.Before(time.Now()) {
		log.Info("Captcha has expired, retrying", "message_id", messageId, "user_id", userId, "chat_id", chatId)
		ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "The captcha has expired",
			ShowAlert: true,
		})
		return ButtonCaptcha(b, captchaState.CaptchaChatId, chatId, userId)
	}

	ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text:      "Thank you for your confirmation! Your account has been approved.",
		ShowAlert: false,
	})
	return ApproveUser(b, &captchaState)
}
