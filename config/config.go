package config

import (
	"time"

	"github.com/spf13/viper"
)

func GetCacheDuration() time.Duration {
	return time.Duration(viper.GetInt("cache.duration")) * time.Minute
}
