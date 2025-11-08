package emojis

import (
	"os"
	"strings"
)

func loadEmojis() (emojis []string, err error) {
	emojis = make([]string, 0)

	files, err := os.ReadDir("./imgs")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := strings.Split(file.Name(), ".")[0]

		emojis = append(emojis, name)
	}

	return
}
