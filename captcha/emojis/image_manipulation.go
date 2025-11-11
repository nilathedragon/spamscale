package emojis

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"

	"github.com/Ostsol/gradient"
	"github.com/nilathedragon/spamscale/resources"
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
	emojiImageData, err := resources.Read(path)
	if err != nil {
		return
	}

	emojiImage, _, err := image.Decode(bytes.NewReader(emojiImageData))
	if err != nil {
		return
	}

	draw.Draw(baseImage, image.Rect(point.X, point.Y, point.X+200, point.Y+200), emojiImage, image.Point{X: 0, Y: 0}, draw.Over)
	return
}
