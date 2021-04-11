package main

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

type ImageCreator struct {
	halfSize int
	seed     int64
	img      *image.NRGBA
	flags    ImageGenFlags
}

type ImageGenFlags struct {
	patterLineWidth, patternModifier float64
	reverseForm, reverseAlpha, reverseRows, colorCombination,
	invert, saturation, contrast, rotate, tunnel, pattern int
}

func NewImageCreator(halfSize int) ImageCreator {
	imageCreator := ImageCreator{halfSize: halfSize}
	imageCreator.randomize()
	return imageCreator
}

func (imageCreator *ImageCreator) randomize() {
	imageCreator.seed = time.Now().UTC().UnixNano()
	rand.Seed(imageCreator.seed)

	imageCreator.flags = ImageGenFlags{
		reverseForm:      rand.Intn(2),
		reverseAlpha:     rand.Intn(2),
		reverseRows:      rand.Intn(2),
		colorCombination: rand.Intn(4),
		invert:           rand.Intn(2),
		saturation:       rand.Intn(30),
		contrast:         rand.Intn(30),
		rotate:           rand.Intn(2),
		tunnel:           rand.Intn(2),
		pattern:          rand.Intn(2),
		patterLineWidth:  float64(rand.Intn(5) + 2),
		patternModifier:  float64(rand.Intn(6) + 5),
	}

	//log.Printf("%+v\n", imageCreator)
}

func (imageCreator *ImageCreator) SavePNG(fileName string) {
	file, err := os.Create(fileName + ".png")
	if err != nil {
		panic(err)
	}

	encoder := &png.Encoder{CompressionLevel: png.NoCompression}
	err = encoder.Encode(file, imageCreator.img)

	if err != nil {
		panic(err)
	}
}

func (imageCreator *ImageCreator) GenImage() {
	size := imageCreator.halfSize * 2
	rect := image.Rect(0, 0, size, size)
	img := image.NewNRGBA(rect)

	red, green, blue := GetRandomRGB()

	pixelSetter := PixelSetter{
		referenceSize: 512,
		halfSize:      imageCreator.halfSize,
		flags:         imageCreator.flags,
		red:           red,
		green:         green,
		blue:          blue,
		img:           img,
	}

	for x := 0; x <= imageCreator.halfSize; x++ {
		pixelSetter.GenLineColor()
		for y := 0; y <= imageCreator.halfSize; y++ {
			if imageCreator.flags.reverseRows == 0 {
				pixelSetter.SetPixels(x, y)
			} else {
				pixelSetter.SetPixels(y, x)
			}
		}
	}

	brightnessDiff := float64(25)
	if imageCreator.flags.invert == 1 {
		img = imaging.Invert(img)
		img = imaging.AdjustBrightness(img, -brightnessDiff)
	} else {
		img = imaging.AdjustBrightness(img, brightnessDiff)
	}

	img = imaging.AdjustSaturation(img, float64(imageCreator.flags.saturation))
	img = imaging.AdjustContrast(img, float64(imageCreator.flags.contrast))

	cropSize := int(float32(imageCreator.halfSize) * 1.4)

	if imageCreator.flags.rotate == 1 {
		img = imaging.Rotate(img, 45, color.Black)

		img = imaging.CropCenter(img, cropSize, cropSize)
		img = imaging.Resize(img, size, size, imaging.Lanczos)
	}

	if imageCreator.flags.rotate == 1 && imageCreator.flags.tunnel == 1 {
		quarterSize := imageCreator.halfSize / 2
		for i := 0; i < 5; i++ {
			cropped := imaging.Resize(img, imageCreator.halfSize, imageCreator.halfSize, imaging.Lanczos)
			img = imaging.Overlay(img, cropped, image.Pt(quarterSize, quarterSize), 1.0)
		}
	}

	imageCreator.img = img
}
