package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
)

func generateRandomGradient(w, h int) (*bytes.Buffer, error) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	color1 := getRandomColor()
	color2 := getRandomColor()

	// Random fill direction
	if rand.Intn(2) == 1 {
		fillGradientLeftToRight(img, color1, color2)
	} else {
		fillGradientTopToBottom(img, color1, color2)
	}

	var body bytes.Buffer

	err := png.Encode(&body, img)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func getRandomColor() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
}

func fillGradientLeftToRight(img *image.RGBA, startColor, endColor color.RGBA) {
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		t := float64(x-bounds.Min.X) / float64(bounds.Max.X-bounds.Min.X)
		r := uint8(float64(startColor.R)*(1-t) + float64(endColor.R)*t)
		g := uint8(float64(startColor.G)*(1-t) + float64(endColor.G)*t)
		b := uint8(float64(startColor.B)*(1-t) + float64(endColor.B)*t)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			color := color.RGBA{r, g, b, 255}
			img.Set(x, y, color)
		}
	}
}

func fillGradientTopToBottom(img *image.RGBA, startColor, endColor color.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		t := float64(y-bounds.Min.Y) / float64(bounds.Max.Y-bounds.Min.Y)
		r := uint8(float64(startColor.R)*(1-t) + float64(endColor.R)*t)
		g := uint8(float64(startColor.G)*(1-t) + float64(endColor.G)*t)
		b := uint8(float64(startColor.B)*(1-t) + float64(endColor.B)*t)
		color := color.RGBA{r, g, b, 255}
		draw.Draw(img, image.Rect(bounds.Min.X, y, bounds.Max.X, y+1), &image.Uniform{color}, image.Point{bounds.Min.X, y}, draw.Src)
	}
}
