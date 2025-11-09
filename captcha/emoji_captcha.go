package captcha

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/png"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/infinytum/injector"
	"github.com/nilathedragon/spamscale/captcha/emojis"
	"github.com/nilathedragon/spamscale/db/model"
	"gorm.io/gorm"
)

const EmojiCaptchaConfirmCallback = "emoji_captcha_confirm"
const errorsBeforeCancelValidation = 3

type EmojiValidationPayload struct {
	DisplayEmojis []string
	ValidEmojis   []string
}

type EmojiValidationState struct {
	Errors                int
	PickedCorrectEmojis   []string
	PickedIncorrectEmojis []string
}

func EmojiCaptcha(b *gotgbot.Bot, captchaChatId int64, chatID int64, userID int64) error {
	correctEmojis, img, err := emojis.GenerateCaptchaImage()
	if err != nil {
		return err
	}

	additionalEmojis, err := emojis.RandomEmojisExcept(correctEmojis, 8)
	if err != nil {
		return err
	}

	emojisToDisplay := append(additionalEmojis, correctEmojis...)

	for i := range emojisToDisplay {
		j := rand.Intn(i + 1)
		emojisToDisplay[i], emojisToDisplay[j] = emojisToDisplay[j], emojisToDisplay[i]
	}

	rows := make([][]gotgbot.InlineKeyboardButton, 6)

	for y := 0; y <= 5; y++ {
		row := make([]gotgbot.InlineKeyboardButton, 2)
		for x := 0; x <= 1; x++ {
			emoji := emojisToDisplay[2*y+x]
			row[x] = gotgbot.InlineKeyboardButton{
				Text:         emoji,
				CallbackData: EmojiCaptchaConfirmCallback + ":" + emoji,
			}
		}
		rows[y] = row
	}

	keyboardMarkup := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	image := bytes.Buffer{}
	err = png.Encode(&image, img)
	if err != nil {
		return err
	}

	message, err := b.SendPhoto(captchaChatId, gotgbot.InputFileByReader("captcha.png", &image), &gotgbot.SendPhotoOpts{
		Caption:     "Please confirm you are not a robot by clicking the corresponding emojis below",
		ReplyMarkup: keyboardMarkup,
	})

	if err != nil {
		return err
	}

	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		return err
	}

	validationPayload, err := json.Marshal(&EmojiValidationPayload{
		DisplayEmojis: emojisToDisplay,
		ValidEmojis:   correctEmojis,
	})

	if err != nil {
		return err
	}

	return db.Create(&model.CaptchaState{
		CaptchaMessageId:  message.MessageId,
		CaptchaChatId:     captchaChatId,
		ChatID:            chatID,
		UserID:            userID,
		ExpiresAt:         time.Now().Add(time.Minute * 10),
		ValidationPayload: string(validationPayload),
	}).Error
}

func EmojiCaptchaCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	messageId := ctx.CallbackQuery.Message.GetMessageId()
	userId := ctx.CallbackQuery.From.Id
	chatId := ctx.CallbackQuery.Message.GetChat().Id
	data, _ := strings.CutPrefix(ctx.CallbackQuery.Data, EmojiCaptchaConfirmCallback+":")

	log.Info("Emoji captcha callback received", "message_id", messageId, "user_id", userId, "chat_id", chatId)

	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		log.Error("An error occurred while injecting the database", "error", err)
		_, _ = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
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
		_, _ = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "An error occurred while processing the captcha",
			ShowAlert: true,
		})
		return err
	}

	// Retry the captcha if it has expired, but the user somehow managed to catch it before the cleanup cronjob ran
	if captchaState.ExpiresAt.Before(time.Now()) {
		log.Info("Captcha has expired, retrying", "message_id", messageId, "user_id", userId, "chat_id", chatId)
		_, _ = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "The captcha has expired",
			ShowAlert: true,
		})
		return EmojiCaptcha(b, captchaState.CaptchaChatId, chatId, userId)
	}

	var currentState EmojiValidationState
	if captchaState.ValidationStatus != "" {
		err = json.Unmarshal([]byte(captchaState.ValidationStatus), &currentState)
		if err != nil {
			return err
		}
	}
	var validationPayload EmojiValidationPayload
	err = json.Unmarshal([]byte(captchaState.ValidationPayload), &validationPayload)
	if err != nil {
		return err
	}

	if slices.Contains(validationPayload.ValidEmojis, data) {
		if slices.Contains(currentState.PickedCorrectEmojis, data) {
			return nil
		}

		currentState.PickedCorrectEmojis = append(currentState.PickedCorrectEmojis, data)
		validationStatus, err := json.Marshal(currentState)
		if err != nil {
			return err
		}

		captchaState.ValidationStatus = string(validationStatus)

		db.Save(captchaState)

		rows := buildKeyboard(validationPayload.DisplayEmojis, currentState)

		keyboardMarkup := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: rows,
		}
		_, _, _ = ctx.CallbackQuery.Message.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
			ChatId:      chatId,
			MessageId:   messageId,
			ReplyMarkup: keyboardMarkup,
		})

		if len(currentState.PickedCorrectEmojis) == len(validationPayload.ValidEmojis) {
			_, _ = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      "Thank you for your confirmation! Your account has been approved.",
				ShowAlert: false,
			})
			return ApproveUser(b, &captchaState)
		}
	} else {
		if slices.Contains(currentState.PickedIncorrectEmojis, data) {
			return nil
		}

		currentState.PickedIncorrectEmojis = append(currentState.PickedIncorrectEmojis, data)
		currentState.Errors++
		validationStatus, err := json.Marshal(currentState)
		if err != nil {
			return err
		}

		captchaState.ValidationStatus = string(validationStatus)
		db.Save(captchaState)

		if currentState.Errors >= errorsBeforeCancelValidation {
			message, err := b.SendMessage(chatId, "Sorry, but you failed the captcha. This message wil delete in 10 seconds.", &gotgbot.SendMessageOpts{})
			if err != nil {
				return err
			}

			time.AfterFunc(10*time.Second, func() {
				_, _ = message.Delete(b, &gotgbot.DeleteMessageOpts{})
			})

			err = RejectUser(b, &captchaState)
			if err != nil {
				return err
			}

			return nil
		}

		rows := buildKeyboard(validationPayload.DisplayEmojis, currentState)

		keyboardMarkup := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: rows,
		}
		_, _, _ = ctx.CallbackQuery.Message.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
			ChatId:      chatId,
			MessageId:   messageId,
			ReplyMarkup: keyboardMarkup,
		})
	}

	return nil
}

func buildKeyboard(emojisToDisplay []string, currentState EmojiValidationState) (rows [][]gotgbot.InlineKeyboardButton) {
	rows = make([][]gotgbot.InlineKeyboardButton, 6)

	for y := 0; y <= 5; y++ {
		row := make([]gotgbot.InlineKeyboardButton, 2)
		for x := 0; x <= 1; x++ {
			emoji := emojisToDisplay[2*y+x]
			if slices.Contains(currentState.PickedCorrectEmojis, emoji) {
				emoji = "✅"
			}
			if slices.Contains(currentState.PickedIncorrectEmojis, emoji) {
				emoji = "❌"
			}

			row[x] = gotgbot.InlineKeyboardButton{
				Text:         emoji,
				CallbackData: EmojiCaptchaConfirmCallback + ":" + emoji,
			}
		}
		rows[y] = row
	}
	return
}
