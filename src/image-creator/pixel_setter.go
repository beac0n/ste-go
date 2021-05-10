package image_creator

import (
	"image"
	"image/color"
	"math"
)

type PixelSetter struct {
	referenceSize, halfSize                                  int
	minAlpha, red, green, blue, redLine, greenLine, blueLine uint8
	img                                                      *image.NRGBA
	flags                                                    ImageGenFlags
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

	if ps.flags.rotate == 1 && ps.flags.pattern == 1 && (ps.calcPatternMatch(x, y) || ps.calcPatternMatch(y, x)) {
		topLeft = GetRandomPixel()
		bottomLeft = GetRandomPixel()
		topRight = GetRandomPixel()
		bottomRight = GetRandomPixel()
	}

	ps.img.SetNRGBA(x, y, topLeft)
	ps.img.SetNRGBA(2*ps.halfSize-x, y, bottomLeft)
	ps.img.SetNRGBA(x, 2*ps.halfSize-y, topRight)
	ps.img.SetNRGBA(2*ps.halfSize-x, 2*ps.halfSize-y, bottomRight)
}

func (ps *PixelSetter) calcPatternMatch(x int, y int) bool {
	xRev := float64(ps.halfSize - x)
	yRev := float64(ps.halfSize - y)

	lineWidth := ps.flags.patterLineWidth * float64(ps.halfSize/ps.referenceSize)
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
	if ps.flags.reverseForm == 1 {
		y *= -1
	}

	alphaModifier := float64(ps.referenceSize) / float64(ps.halfSize)
	alpha := uint8(float64(x+y) * alphaModifier)

	if ps.flags.reverseAlpha == 1 {
		alpha = 255 - alpha
	}

	if alpha < ps.minAlpha {
		alpha = ps.minAlpha
	}

	return alpha
}
