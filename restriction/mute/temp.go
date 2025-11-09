package mute

import (
	"database/sql"
	"regexp"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/infinytum/injector"
	"github.com/nilathedragon/spamscale/db/model"
	"gorm.io/gorm"
)

func TemporaryMute(b *gotgbot.Bot, chatId, userId int64, duration string, ctx *ext.Context) error {
	db, err := injector.Inject[*gorm.DB]()
	if err != nil {
		return err
	}

	regex, err := regexp.Compile(`(\d*)(min|h|d|m|y)`)
	if err != nil {
		return err
	}

	matches := regex.FindStringSubmatch(duration)

	t, err := strconv.Atoi(matches[1])
	t64 := int64(t)
	if err != nil {
		return err
	}
	multiplier := matches[2]

	var deltaTime time.Duration

	switch multiplier {
	case "min":
		deltaTime = time.Duration(t64 * int64(time.Minute))
	case "h":
		deltaTime = time.Duration(t64 * int64(time.Hour))
	case "d":
		deltaTime = time.Duration(t64 * int64(time.Hour*24))
	case "m":
		deltaTime = time.Duration(t64 * int64(time.Hour*24*30))
	case "y":
		deltaTime = time.Duration(t64 * int64(time.Hour*24*365))
	}

	expiration := time.Now().Add(deltaTime)

	if _, err := ctx.Message.Chat.RestrictMember(
		b,
		userId,
		gotgbot.ChatPermissions{},
		&gotgbot.RestrictChatMemberOpts{
			UntilDate: expiration.Unix(),
		},
	); err != nil {
		return err
	}

	if err := db.Save(&model.ChatRestriction{
		ChatID:          chatId,
		UserID:          userId,
		RestrictionType: model.ChatRestrictionTypeMute,
		ExpiresAt: sql.NullInt64{
			Int64: expiration.Unix(),
			Valid: true,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}
