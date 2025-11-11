package emojis

import (
	"image"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

var emojis []string

func GenerateCaptchaImage() (chosenEmojis []string, img *image.RGBA, err error) {
	if len(emojis) == 0 {
		emojis, err = loadEmojis()
		if err != nil {
			return nil, nil, err
		}
	}

	img = createBaseImage()
	xOffset := 140
	alreadyPicked := make([]int, 0, 4)

	for i := 0; i <= 3; i++ {
		var randomIndex int
		for keepGoing := true; keepGoing; keepGoing = slices.Contains(alreadyPicked, randomIndex) {
			randomIndex = int(rand.Float64() * float64(len(emojis)))
		}

		alreadyPicked = append(alreadyPicked, randomIndex)

		emoji := strings.Join(emojiToUTF8(emojis[randomIndex]), "_")

		err := addImage(img, image.Point{X: xOffset, Y: 200}, viper.GetString("resources-dir")+"/emojis/"+emoji+".png")
		if err != nil {
			return nil, nil, err
		}
		xOffset += 220
	}

	chosenEmojis = make([]string, 4)

	for i, index := range alreadyPicked {
		chosenEmojis[i] = emojis[index]
	}

	return chosenEmojis, img, nil
}

func RandomEmojiExcept(except []string) (emoji string, err error) {
	emojis, err = loadEmojis()
	if err != nil {
		return "", err
	}

	exceptIndices := make([]int, len(except))
	for i, emoji := range except {
		exceptIndices[i] = slices.Index(emojis, emoji)
	}

	var randomIndex int
	for keepGoing := true; keepGoing; keepGoing = slices.Contains(exceptIndices, randomIndex) {
		randomIndex = int(rand.Float64() * float64(len(emojis)))
	}

	emoji = emojis[randomIndex]
	return
}

func RandomEmoji() (emoji string, err error) {
	return RandomEmojiExcept([]string{})
}

func RandomEmojisExcept(except []string, count int) (emojis []string, err error) {
	emojis = make([]string, count)

	for i := range emojis {
		emoji, err := RandomEmojiExcept(except)
		if err != nil {
			return nil, err
		}

		emojis[i] = emoji
	}
	return
}

func RandomEmojis(count int) (emojis []string, err error) {
	return RandomEmojisExcept([]string{}, count)
}
