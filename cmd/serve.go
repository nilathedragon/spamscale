package cmd

import (
	"errors"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/go-mojito/mojito/log"
	"github.com/nilathedragon/spamscale/captcha"
	"github.com/nilathedragon/spamscale/handler"
	"github.com/nilathedragon/spamscale/sidecar"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userCommandScope = []gotgbot.BotCommand{
	{
		Command:     "report",
		Description: "Notify the moderators",
	},
	{
		Command:     "boop",
		Description: "Boop the derg",
	},
}

var adminCommandScope = append([]gotgbot.BotCommand{
	{
		Command:     "ban",
		Description: "Ban a user",
	},
	{
		Command:     "setcaptchatype",
		Description: "Set the type of captcha to use",
	},
	{
		Command:     "setfastblocklistenabled",
		Description: "Set whether to use the fast blocklist",
	},
	{
		Command:     "tmute",
		Description: "Temporarily mute a user. Usage: /tmute <@username> <duration>",
	},
}, userCommandScope...)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the SpamScale telegram bot",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Starting SpamScale telegram bot...")

		botToken := viper.GetString("bot-token")
		if !viper.IsSet("bot-token") || botToken == "" {
			return errors.New("bot-token is not set")
		}

		bot, err := gotgbot.NewBot(botToken, nil)
		if err != nil {
			return err
		}

		dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Error("an error occurred while handling update", "error", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		})
		updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{})

		dispatcher.AddHandler(handlers.NewMyChatMember(handler.SetupHandlerFilter, handler.SetupHandler))
		dispatcher.AddHandler(handlers.NewChatJoinRequest(handler.JoinRequestHandlerFilter, handler.JoinRequestHandler))
		dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix(captcha.ButtonCaptchaConfirmCallback), captcha.ButtonCaptchaCallback))
		dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix(captcha.EmojiCaptchaConfirmCallback), captcha.EmojiCaptchaCallback))
		dispatcher.AddHandler(handlers.NewCommand("ban", handler.CommandBanHandler))
		dispatcher.AddHandler(handlers.NewCommand("boop", handler.CommandBoopHandler))
		dispatcher.AddHandler(handlers.NewCommand("captcha", handler.CommandCaptchaHandler))
		dispatcher.AddHandler(handlers.NewCommand("tmute", handler.CommandTMuteHandler))
		dispatcher.AddHandler(handlers.NewCommand("pet", handler.CommandPetHandler))
		dispatcher.AddHandler(handlers.NewCommand("report", handler.CommandReportHandler))
		dispatcher.AddHandler(handlers.NewCommand("setcaptchatype", handler.CommandSetCaptchaTypeHandler))
		dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix(handler.CommandSetCaptchaTypeCallback), handler.CommandSetCaptchaTypeHandlerCallback))
		dispatcher.AddHandler(handlers.NewCommand("setfastblocklistenabled", handler.CommandSetFastBlocklistEnabledHandler))
		dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix(handler.CommandSetFastBlocklistEnabledCallback), handler.CommandSetFastBlocklistEnabledHandlerCallback))

		if _, err := bot.SetMyCommands(
			userCommandScope,
			&gotgbot.SetMyCommandsOpts{
				Scope: gotgbot.BotCommandScopeDefault{},
			},
		); err != nil {
			return err
		}

		if _, err := bot.SetMyCommands(
			adminCommandScope,
			&gotgbot.SetMyCommandsOpts{
				Scope: gotgbot.BotCommandScopeAllChatAdministrators{},
			},
		); err != nil {
			return err
		}

		err = updater.StartPolling(bot, &ext.PollingOpts{
			DropPendingUpdates: false,
			GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
				Timeout: 10,
				RequestOpts: &gotgbot.RequestOpts{
					Timeout: time.Second * 11,
				},
			},
		})
		if err != nil {
			return err
		}

		go sidecar.ExpireCaptchas(bot)

		log.Info("SpamScale telegram bot started, press Ctrl+C to stop")
		updater.Idle()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
