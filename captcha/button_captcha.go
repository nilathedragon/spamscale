package captcha

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/db"
	"github.com/nilathedragon/spamscale/db/model"
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

	return db.CaptchaState.Save(&model.CaptchaState{
		CaptchaMessageId: message.MessageId,
		CaptchaChatId:    captchaChatId,
		ChatID:           chatID,
		UserID:           userID,
		ExpiresAt:        time.Now().Add(time.Minute * 10),
	})
}

func ButtonCaptchaCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	captchaState, err := db.CaptchaState.Get(ctx.CallbackQuery.Message.GetChat().Id, ctx.CallbackQuery.Message.GetMessageId(), ctx.CallbackQuery.From.Id)
	if err != nil {
		ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "An error occurred while processing the captcha, if the error persists, please try again or contact the administrator.",
			ShowAlert: true,
		})
		return err
	}
	log.Info("Button captcha callback received", "message_id", captchaState.CaptchaMessageId, "user_id", captchaState.UserID, "chat_id", captchaState.ChatID)

	// The button captcha does not provide any additional validation, so we can approve the user immediately
	ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text:      "Thank you for your confirmation! Your join request has been approved.",
		ShowAlert: false,
	})

	return ApproveUser(b, &captchaState)
}
