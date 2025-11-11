package emojis

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func loadEmojis() (emojis []string, err error) {
	emojis = make([]string, 0)

	files, err := os.ReadDir(viper.GetString("resources-dir") + "/emojis")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := strings.Split(file.Name(), ".")[0]

		emojis = append(emojis, name)
	}

	return
}
