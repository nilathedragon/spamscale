package emojis

import (
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/Ostsol/gradient"
)

func createBaseImage() (img *image.RGBA) {
	img = image.NewRGBA(image.Rect(0, 0, 1140, 600))
	gradient.DrawLinear(
		img,
		0,
		0,
		1,
		1,
		[]gradient.Stop{
			{
				X:   0,
				Col: color.RGBA{R: 0x92, G: 0x0D, B: 0xE9, A: 0xFF},
			},
			{
				X:   1,
				Col: color.RGBA{R: 0x15, G: 0x36, B: 0xF1, A: 0xFF},
			},
		},
	)

	return
}

func addImage(baseImage draw.Image, point image.Point, path string) (err error) {
	emojiImageFile, err := os.OpenFile(path, os.O_RDONLY, 0o006)
	if err != nil {
		return
	}
	defer func(emojiImageFile *os.File) {
		_ = emojiImageFile.Close()
	}(emojiImageFile)

	emojiImage, _, err := image.Decode(emojiImageFile)
	if err != nil {
		return
	}
	defer func(emojiImageFile *os.File) {
		_ = emojiImageFile.Close()
	}(emojiImageFile)

	draw.Draw(baseImage, image.Rect(point.X, point.Y, point.X+200, point.Y+200), emojiImage, image.Point{X: 0, Y: 0}, draw.Over)
	return
}
