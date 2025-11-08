package model

import "database/sql"

type ChatRestrictionType string

const (
	ChatRestrictionTypeBan  ChatRestrictionType = "ban"
	ChatRestrictionTypeMute ChatRestrictionType = "mute"
	ChatRestrictionTypeWarn ChatRestrictionType = "warn"
)

// Represents a moderation action taken in a chat by the bot
type ChatRestriction struct {
	// The chat in which the action was taken
	ChatID int64 `gorm:"primarykey"`
	// The user against which an action was taken
	UserID int64 `gorm:"primarykey"`
	// Action taken against the user
	RestrictionType ChatRestrictionType `gorm:"not null"`
	// Time at which this restriction will be lifted, if at all
	ExpiresAt sql.NullInt64
}
