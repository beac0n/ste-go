package main

import (
	"image"
	"image/color"
	"math"
)

type PixelSetter struct {
	halfSize                                       int
	red, green, blue, redLine, greenLine, blueLine uint8
	img                                            *image.NRGBA
	flags                                          ImageGenFlags
}

func (ps *PixelSetter) GenLineColor() {
	redLine, greenLine, blueLine := GetRandomRGB()
	ps.redLine = redLine
	ps.greenLine = greenLine
	ps.blueLine = blueLine
}

func (ps *PixelSetter) SetPixels(x, y int) {
	baseColor := ps.getBaseColor(x, y)

	topLeft := ps.createColorRGBA(baseColor)
	bottomLeft := ps.createColorRGBA(baseColor)
	topRight := ps.createColorRGBA(baseColor)
	bottomRight := ps.createColorRGBA(baseColor)

	if ps.flags.rotate == 1 && (ps.calcPatternMatch(x, y) || ps.calcPatternMatch(y, x)) {
		patternColor := color.NRGBA{
			A: 255,
			R: ps.blue & GetRandomColorValue(),
			G: ps.red & GetRandomColorValue(),
			B: ps.green & GetRandomColorValue(),
		}
		topLeft = patternColor
		bottomLeft = patternColor
		topRight = patternColor
		bottomRight = patternColor
	}

	ps.img.SetNRGBA(x, y, topLeft)
	ps.img.SetNRGBA(2*ps.halfSize-x, y, bottomLeft)
	ps.img.SetNRGBA(x, 2*ps.halfSize-y, topRight)
	ps.img.SetNRGBA(2*ps.halfSize-x, 2*ps.halfSize-y, bottomRight)
}

func (ps *PixelSetter) calcPatternMatch(x int, y int) bool {
	xRev := float64(ps.halfSize - x)
	yRev := float64(ps.halfSize - y)

	lineWidth := ps.flags.patterLineWidth
	zModifier := ps.flags.patternModifier
	z0 := float64(ps.halfSize) / zModifier
	z1 := float64(ps.halfSize) * 0.9

	expectedX0 := math.Sin((xRev/z0)-math.Pi) * z1
	expectedX1 := math.Sin(xRev/z0) * z1

	isOnLine0 := expectedX0 >= yRev-lineWidth && expectedX0 <= yRev+lineWidth
	isOnLine1 := expectedX1 >= yRev-lineWidth && expectedX1 <= yRev+lineWidth

	return isOnLine0 || isOnLine1
}

func (ps *PixelSetter) createColorRGBA(baseColor uint8) color.NRGBA {
	baseColor = baseColor / 2

	rgb := [3]uint8{
		ps.red & ps.redLine & GetRandomColorValue(),
		ps.green & ps.greenLine & GetRandomColorValue(),
		ps.blue & ps.blueLine & GetRandomColorValue(),
	}

	for i := 0; i < 3; i++ {
		if i != ps.flags.colorCombination {
			rgb[i] = baseColor ^ rgb[i]
		} else {
			rgb[i] = baseColor & rgb[i]
		}
	}

	return color.NRGBA{A: 255, R: rgb[0], G: rgb[1], B: rgb[2]}
}

func (ps *PixelSetter) getBaseColor(x int, y int) uint8 {
	var alpha uint8

	if ps.flags.reverseForm == 1 {
		alpha = uint8(x - y)
	} else {
		alpha = uint8(x + y)
	}

	if ps.flags.reverseAlpha == 1 {
		alpha = 255 - alpha
	}

	minAlpha := uint8(50)
	if alpha < minAlpha {
		return minAlpha
	} else {
		return alpha
	}
}
