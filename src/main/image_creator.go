package main

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strconv"
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
		rotate:           1, //rand.Intn(2),
		tunnel:           1, //rand.Intn(2),
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

func (imageCreator *ImageCreator) Encode(data []byte) error {
	if imageCreator.img == nil {
		return errors.New("no image to encode data in")
	}

	fullSize := imageCreator.halfSize * 4
	maximumDataLength := (fullSize*fullSize - 32) / 8

	if len(data) >= maximumDataLength {
		return errors.New("length of provided data is bigger than the maximum length: " + strconv.Itoa(maximumDataLength))
	}

	for i, dataLengthBitChar := range fmt.Sprintf("%032b", len(data)*8) {
		dataLengthBit := byte(int(dataLengthBitChar - '0'))
		imageCreator.encodeBitInByte(i, dataLengthBit)
	}

	bitIndex := 32
	for _, dataByte := range data {
		for j := 0; j < 8; j++ {
			bitmask := byte(math.Exp2(float64(j)))
			maskedDataByte := dataByte & bitmask
			dataBit := byte(0)
			if maskedDataByte > 0 {
				dataBit = byte(1)
			}

			// make sure to skip A values
			imageCreator.encodeBitInByte(bitIndex, dataBit)
			bitIndex++
		}
	}

	return nil
}

func (imageCreator *ImageCreator) Decode(img *image.NRGBA) ([]byte, error) {
	if img == nil {
		return nil, errors.New("no image provided")
	}

	dataInBitsLength := 0
	for i := 0; i < 32; i++ {
		encodedBit := int(img.Pix[i] & 1)
		dataInBitsLength += encodedBit * int(math.Exp2(float64(31-i)))
	}

	data := make([]byte, dataInBitsLength/8)

	dataBitIndex := dataInBitsLength + 32
	for i := 32; i < dataBitIndex; i += 8 {
		encodedByte := byte(0)
		for j := 0; j < 8; j++ {
			encodedBit := img.Pix[i+j] & 1
			encodedByte += encodedBit * byte(math.Exp2(float64(j)))
		}
		data[(i-32)/8] = encodedByte
	}

	return data, nil
}

func (imageCreator *ImageCreator) encodeBitInByte(byteIndex int, bit byte) {
	if bit == 0 {
		imageCreator.img.Pix[byteIndex] = imageCreator.img.Pix[byteIndex] & 254
	} else {
		imageCreator.img.Pix[byteIndex] = imageCreator.img.Pix[byteIndex] | 1
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
		tunnelPosition := imageCreator.halfSize / 2
		for i := 0; i < 5; i++ {
			cropped := imaging.Resize(img, imageCreator.halfSize, imageCreator.halfSize, imaging.Lanczos)
			img = imaging.Overlay(img, cropped, image.Pt(tunnelPosition, tunnelPosition), 1.0)
		}
	}

	imageCreator.img = img
}
