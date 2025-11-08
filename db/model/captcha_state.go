package model

import "time"

type CaptchaState struct {
	// The message ID of the message that contains the captcha trigger
	CaptchaMessageId int64 `gorm:"primarykey"`
	// The chat id of the chat in which the captcha message was ssent
	CaptchaChatId int64 `gorm:"primarykey"`
	// The user id of the user trying to join
	UserID int64 `gorm:"primarykey"`

	// The chat id of the chat in which the user is trying to join
	ChatID int64 `gorm:"not null"`

	// The payload used to validate the verification of the user
	ValidationPayload string
	// The current state of the validation
	ValidationStatus string

	// The time at which the captcha state will expire
	ExpiresAt time.Time `gorm:"not null"`
}
