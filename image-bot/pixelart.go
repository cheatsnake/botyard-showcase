package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
)

func newPixelart(width, height, boxSize int) (*bytes.Buffer, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x += boxSize {
		for y := 0; y < height; y += boxSize {
			boxColor := randomColor()
			drawBox(img, x, y, boxSize, boxColor)
		}
	}

	var body bytes.Buffer
	err := png.Encode(&body, img)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func randomColor() color.Color {
	r := uint8(rand.Intn(256))
	g := uint8(rand.Intn(256))
	b := uint8(rand.Intn(256))
	return color.RGBA{r, g, b, 255}
}

func drawBox(img *image.RGBA, x, y, boxSize int, c color.Color) {
	for i := x; i < x+boxSize; i++ {
		for j := y; j < y+boxSize; j++ {
			img.Set(i, j, c)
		}
	}
}
