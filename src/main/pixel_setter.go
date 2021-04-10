package main

import (
	"image"
	"image/color"
)

type PixelSetter struct {
	size, reverseForm, reverseAlpha, colorCombination int
	red, green, blue, redLine, greenLine, blueLine    uint8
	img                                               *image.NRGBA
}

func (ps *PixelSetter) GenLineColor() {
	redLine, greenLine, blueLine := GetRandomRGB()
	ps.redLine = redLine
	ps.greenLine = greenLine
	ps.blueLine = blueLine
}

func (ps *PixelSetter) SetPixels(x, y int) {
	baseColor := ps.getBaseColor(x, y)

	ps.img.SetNRGBA(x, y, ps.createColorRGBA(baseColor))
	ps.img.SetNRGBA(2*ps.size-x, y, ps.createColorRGBA(baseColor))
	ps.img.SetNRGBA(x, 2*ps.size-y, ps.createColorRGBA(baseColor))
	ps.img.SetNRGBA(2*ps.size-x, 2*ps.size-y, ps.createColorRGBA(baseColor))
}

func (ps *PixelSetter) createColorRGBA(baseColor uint8) color.NRGBA {
	baseColor = baseColor / 2

	rgb := [3]uint8{
		ps.red & ps.redLine & GetRandomColorValue(),
		ps.green & ps.greenLine & GetRandomColorValue(),
		ps.blue & ps.blueLine & GetRandomColorValue(),
	}

	for i := 0; i < 3; i++ {
		if i != ps.colorCombination {
			rgb[i] = baseColor ^ rgb[i]
		} else {
			rgb[i] = baseColor & rgb[i]
		}
	}

	return color.NRGBA{A: 255, R: rgb[0], G: rgb[1], B: rgb[2]}
}

func (ps *PixelSetter) getBaseColor(x int, y int) uint8 {
	var alpha uint8

	if ps.reverseForm == 1 {
		alpha = uint8(x - y)
	} else {
		alpha = uint8(x + y)
	}

	if ps.reverseAlpha == 1 {
		alpha = 255 - alpha
	}

	minAlpha := uint8(50)
	if alpha < minAlpha {
		return minAlpha
	} else {
		return alpha
	}
}
