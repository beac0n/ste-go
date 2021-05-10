package image_creator

import (
	"image/color"
	"math/rand"
)

func GetRandomRGB() (uint8, uint8, uint8) {
	red := GetRandomColorValue()
	green := GetRandomColorValue()
	blue := GetRandomColorValue()
	return red, green, blue
}

func GetRandomColorValue() uint8 {
	return uint8(rand.Intn(256))
}

func GetRandomPixel() color.NRGBA {
	return color.NRGBA{
		A: GetRandomColorValue(),
		R: GetRandomColorValue(),
		G: GetRandomColorValue(),
		B: GetRandomColorValue(),
	}
}
