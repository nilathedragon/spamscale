package emojis

import (
	"os"
	"strconv"
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
		emojis = append(emojis, utf8ToEmoji(strings.Split(name, "_")))

	}

	return
}

// UTF8ToEmoji returns the emoji (in a string) corresponding to the UTF8 codes passed.
//
// The codes are in format XXXX_XXXX where XXXX corresponds to the U+XXXX representation.
func utf8ToEmoji(UTF8 []string) string {
	emoji := make([]rune, 0)

	for _, char := range UTF8 {
		num, _ := strconv.ParseInt(char, 16, 64)
		emoji = append(emoji, rune(num))
	}

	return string(emoji)
}

// EmojiToUTF8 returns the UTF8 codes corresponding to the emoji passed.
//
// The codes are in format XXXX_XXXX where XXXX corresponds to the U+XXXX representation.
func emojiToUTF8(emoji string) []string {
	runes := []rune(emoji)
	codes := make([]string, len(runes))

	for i, r := range runes {
		codes[i] = strconv.FormatInt(int64(r), 16)
	}

	return codes
}
