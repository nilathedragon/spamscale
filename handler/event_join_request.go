package handler

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/captcha"
	"github.com/nilathedragon/spamscale/db"
	"github.com/nilathedragon/spamscale/pkg/thirdparty/fast"
)

func JoinRequestHandlerFilter(cjr *gotgbot.ChatJoinRequest) bool {
	return true
}

type PetData struct {
	Cancel           func()
	MessagesToDelete []int64
}

const (
	PetsTimeout = 30 * time.Second
)

var (
	PetCommand = gotgbot.BotCommand{
		Command:     "pet",
		Description: "Pet the derg",
	}

	petCommandAllowedMap = make(map[int64]*PetData)
)

func petMapHandler(b *gotgbot.Bot, ctx context.Context, chatID int64) {
	petCommandAllowedMapEntry := petCommandAllowedMap[chatID]
	defer func() {
		for _, msgId := range petCommandAllowedMapEntry.MessagesToDelete {
			b.DeleteMessage(chatID, msgId, &gotgbot.DeleteMessageOpts{})
		}
	}()
	select {
	case <-time.After(PetsTimeout):
		commands, err := b.GetMyCommands(&gotgbot.GetMyCommandsOpts{
			Scope: &gotgbot.BotCommandScopeChat{
				ChatId: chatID,
			},
		})

		filteredCommands := make([]gotgbot.BotCommand, 0, len(commands))
		for _, cmd := range commands {
			if cmd.Command != PetCommand.Command {
				filteredCommands = append(filteredCommands, cmd)
			}
		}

		if len(filteredCommands) == 0 {
			_, err = b.DeleteMyCommands(&gotgbot.DeleteMyCommandsOpts{
				Scope: gotgbot.BotCommandScopeChat{
					ChatId: chatID,
				},
			})
		} else {
			if _, err = b.SetMyCommands(filteredCommands, &gotgbot.SetMyCommandsOpts{
				Scope: &gotgbot.BotCommandScopeChat{
					ChatId: chatID,
				},
			}); err != nil {
				return
			}
		}

		commands, err = b.GetMyCommands(&gotgbot.GetMyCommandsOpts{
			Scope: &gotgbot.BotCommandScopeChatAdministrators{
				ChatId: chatID,
			},
		})
		filteredCommands = make([]gotgbot.BotCommand, 0, len(commands))

		for _, cmd := range commands {
			if cmd.Command != PetCommand.Command {
				filteredCommands = append(filteredCommands, cmd)
			}
		}

		if len(filteredCommands) == 0 {
			_, err = b.DeleteMyCommands(&gotgbot.DeleteMyCommandsOpts{
				Scope: gotgbot.BotCommandScopeChatAdministrators{
					ChatId: chatID,
				},
			})
		} else {
			if _, err = b.SetMyCommands(filteredCommands, &gotgbot.SetMyCommandsOpts{
				Scope: &gotgbot.BotCommandScopeChatAdministrators{
					ChatId: chatID,
				},
			}); err != nil {
				return
			}
		}
		fmt.Println("DELETING COMMANDS")
		delete(petCommandAllowedMap, chatID)
	case <-ctx.Done():
		return
	}
}

func JoinRequestHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	logger := log.With("chat_id", ctx.ChatJoinRequest.Chat.Id, "user_id", ctx.ChatJoinRequest.From.Id)
	logger.Info("Chat join request received")

	if fastBlocklistEnabled, err := db.Chat.IsFastBlocklistEnabled(ctx.ChatJoinRequest.Chat.Id); err != nil {
		logger.Error("Failed to check if fast blocklist is enabled", "error", err)
	} else if fastBlocklistEnabled {
		if blocked, err := fast.IsBlocked(ctx.ChatJoinRequest.From.Id); err != nil {
			return err
		} else if blocked {
			logger.Info("User is blocked by fast blocklist, rejecting join request")

			return BlockUser(b, ctx.ChatJoinRequest.Chat.Id, ctx.ChatJoinRequest.From.Id)
		}
	}

	return captcha.TriggerCaptcha(b, ctx.ChatJoinRequest.UserChatId, ctx.ChatJoinRequest.Chat.Id, ctx.ChatJoinRequest.From.Id)
}

func BlockUser(b *gotgbot.Bot, chatID, userID int64) error {
	message, err := b.SendMessage(chatID, "A probable scammer just got blocked!", &gotgbot.SendMessageOpts{})
	if err != nil {
		return err
	}

	if _, exists := petCommandAllowedMap[chatID]; exists {
		petCommandAllowedMap[chatID].Cancel()
		cmdCtx, cancel := context.WithCancel(context.Background())
		petCommandAllowedMap[chatID] = &PetData{
			Cancel:           cancel,
			MessagesToDelete: []int64{message.MessageId},
		}

		go petMapHandler(b, cmdCtx, chatID)
		return nil
	}

	_, _ = b.DeleteMyCommands(&gotgbot.DeleteMyCommandsOpts{
		Scope: gotgbot.BotCommandScopeChat{
			ChatId: chatID,
		},
	})

	_, _ = b.DeleteMyCommands(&gotgbot.DeleteMyCommandsOpts{
		Scope: gotgbot.BotCommandScopeChatAdministrators{
			ChatId: chatID,
		},
	})

	cmds, err := b.GetMyCommands(&gotgbot.GetMyCommandsOpts{
		Scope: gotgbot.BotCommandScopeChat{ChatId: chatID},
	})
	if err != nil {
		return err
	} else if len(cmds) == 0 {
		cmds, _ = b.GetMyCommands(&gotgbot.GetMyCommandsOpts{
			Scope: gotgbot.BotCommandScopeDefault{},
		})
	}
	if !slices.Contains(cmds, PetCommand) {
		cmds = append(cmds, PetCommand)
	}
	_, err = b.SetMyCommands(cmds, &gotgbot.SetMyCommandsOpts{
		Scope: gotgbot.BotCommandScopeChat{ChatId: chatID},
	})
	if err != nil {
		return err
	}

	adminCmds, err := b.GetMyCommands(&gotgbot.GetMyCommandsOpts{
		Scope: gotgbot.BotCommandScopeChatAdministrators{ChatId: chatID},
	})
	if err != nil {
		return err
	} else if len(adminCmds) == 0 {
		adminCmds, _ = b.GetMyCommands(&gotgbot.GetMyCommandsOpts{
			Scope: gotgbot.BotCommandScopeAllChatAdministrators{},
		})

		if len(adminCmds) == 0 {
			adminCmds = slices.Clone(cmds)
		}
	}
	if !slices.Contains(adminCmds, PetCommand) {
		adminCmds = append(adminCmds, PetCommand)
	}
	if _, err = b.SetMyCommands(adminCmds, &gotgbot.SetMyCommandsOpts{
		Scope: gotgbot.BotCommandScopeChatAdministrators{ChatId: chatID},
	}); err != nil {
		return err
	}

	cmdCtx, cancel := context.WithCancel(context.Background())
	petCommandAllowedMap[chatID] = &PetData{
		Cancel: cancel,
		MessagesToDelete: []int64{
			message.MessageId,
		},
	}

	go petMapHandler(b, cmdCtx, chatID)

	if _, err := b.DeclineChatJoinRequest(chatID, userID, &gotgbot.DeclineChatJoinRequestOpts{}); err != nil {
		return err
	}
	return nil
}
