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
	halfSize, reverseForm, reverseAlpha, reverseRows, colorCombination,
	invert, saturation, contrast, rotate, tunnel int
	seed int64
	img  *image.NRGBA
}

func NewImageCreator(size int) ImageCreator {
	imageCreator := ImageCreator{halfSize: size}
	imageCreator.randomize()
	return imageCreator
}

func (imageCreator *ImageCreator) randomize() {
	imageCreator.seed = time.Now().UTC().UnixNano()
	rand.Seed(imageCreator.seed)

	imageCreator.reverseForm = rand.Intn(2)
	imageCreator.reverseAlpha = rand.Intn(2)
	imageCreator.reverseRows = rand.Intn(2)
	imageCreator.colorCombination = rand.Intn(4)
	imageCreator.invert = rand.Intn(2)
	imageCreator.saturation = rand.Intn(30)
	imageCreator.contrast = rand.Intn(30)
	imageCreator.rotate = rand.Intn(2)
	imageCreator.tunnel = rand.Intn(2)

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

func (imageCreator *ImageCreator) GetImage() {
	size := imageCreator.halfSize * 2
	rect := image.Rect(0, 0, size, size)
	img := image.NewNRGBA(rect)

	red, green, blue := GetRandomRGB()

	pixelSetter := PixelSetter{
		size:             imageCreator.halfSize,
		reverseForm:      imageCreator.reverseForm,
		reverseAlpha:     imageCreator.reverseAlpha,
		colorCombination: imageCreator.colorCombination,
		red:              red,
		green:            green,
		blue:             blue,
		img:              img,
	}

	for x := 0; x <= imageCreator.halfSize; x++ {
		pixelSetter.GenLineColor()
		for y := 0; y <= imageCreator.halfSize; y++ {
			if imageCreator.reverseRows == 0 {
				x, y = ReverseValues(x, y)
			}
			pixelSetter.SetPixels(x, y)
		}
	}

	brightnessDiff := float64(25)
	if imageCreator.invert == 1 {
		img = imaging.Invert(img)
		img = imaging.AdjustBrightness(img, -brightnessDiff)
	} else {
		img = imaging.AdjustBrightness(img, brightnessDiff)
	}

	img = imaging.AdjustSaturation(img, float64(imageCreator.saturation))
	img = imaging.AdjustContrast(img, float64(imageCreator.contrast))

	cropSize := int(float32(imageCreator.halfSize) * 1.4)

	if imageCreator.rotate == 1 {
		img = imaging.Rotate(img, 45, color.Black)

		img = imaging.CropCenter(img, cropSize, cropSize)
		img = imaging.Resize(img, size, size, imaging.Lanczos)
	}

	if imageCreator.rotate == 1 && imageCreator.tunnel == 1 {
		quarterSize := imageCreator.halfSize / 2
		for i := 0; i < 5; i++ {
			cropped := imaging.Resize(img, imageCreator.halfSize, imageCreator.halfSize, imaging.Lanczos)
			img = imaging.Overlay(img, cropped, image.Pt(quarterSize, quarterSize), 1.0)
		}
	}

	imageCreator.img = img
}
