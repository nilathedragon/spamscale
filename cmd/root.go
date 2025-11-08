package cmd

import (
	"os"
	"strings"

	"github.com/infinytum/injector"
	"github.com/nilathedragon/spamscale/db/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "spa",
	Short: "SpamScale is a Telegram group chat moderation tool",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Configure Viper to read the config file
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/spamscale")
		viper.SetConfigType("yaml")
		viper.SetConfigName("spamscale")
		viper.AutomaticEnv()
		viper.SetEnvPrefix("SPAMSCALE")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		// Configure the dependency injector to supply a database connection
		injector.DeferredSingleton(func() *gorm.DB {
			db, err := gorm.Open(sqlite.Open("spamscale.db"), &gorm.Config{})
			if err != nil {
				panic(err)
			}

			if err := db.AutoMigrate(&model.CaptchaState{}, &model.Chat{}); err != nil {
				panic(err)
			}
			return db
		})

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("bot-token", "t", "", "Telegram API bot token")
	viper.BindPFlag("bot-token", rootCmd.Flags().Lookup("bot-token"))
}
