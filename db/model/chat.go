package model

import "database/sql"

type CaptchaFailAction string

const (
	CaptchaFailActionBan  CaptchaFailAction = "ban"
	CaptchaFailActionKick CaptchaFailAction = "kick"
	CaptchaFailActionMute CaptchaFailAction = "mute"
)

type CaptchaType string

const (
	CaptchaTypeNone   CaptchaType = "none"
	CaptchaTypeButton CaptchaType = "button"
	CaptchaTypeEmoji  CaptchaType = "emoji"
	CaptchaTypeManual CaptchaType = "manual"
)

// Represents a Telegram chat
type Chat struct {
	ID int64 `gorm:"primarykey"`

	// Join Captcha Settings

	// Whether the join captcha is enabled or not
	CaptchaType CaptchaType `gorm:"not null;default:button"`
	// How much time, in minutes, a user has to solve the join captcha
	CaptchaTimeout int `gorm:"not null;default:10"`
	// What to do when a user fails the join captcha
	CaptchaFailAction CaptchaFailAction `gorm:"not null;default:kick"`
	// How long the user should be banned for, in minutes, after failing the captcha
	CaptchaFailBanDuration sql.NullInt64 `gorm:"default:1440"`
	// How long the user should be muted for, in minutes, after failing the captcha
	CaptchaFailMuteDuration sql.NullInt64 `gorm:"default:60"`

	// Blocklist Settings

	// Whether to ban users present on the Furry Assisted Scam Tracking blocklist
	BlocklistFastEnabled bool `gorm:"not null;default:false"`
}
