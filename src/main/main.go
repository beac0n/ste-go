package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"ste-go/src/image-creator"
	"strings"
)

func main() {
	// TODO: add arguments
	// 	- flags

	dataFilePath := flag.String("data-file-path", "", "path to file to encode in image(s) or path to image with encoded data")
	encode := flag.Bool("encode", true, "encode data (default)")
	decode := flag.Bool("decode", false, "decode data")

	flag.Parse()

	if *dataFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	fileStat, err := os.Stat(*dataFilePath)
	if os.IsNotExist(err) {
		printAndExit("ERROR: "+*dataFilePath+" does not exist", err)
	}

	if *decode {
		decodeData(*dataFilePath)
	} else if *encode {
		fileSize := fileStat.Size()
		imgSize := image_creator.GetImageSizeForByteLength(fileSize)
		encodeData(*dataFilePath, imgSize)
	}
}

func decodeData(dataFilePath string) {
	f, _ := os.Open(dataFilePath)
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		printAndExit("ERROR: could not read file to image:", err)
	}
	data, err := image_creator.Decode(img.(*image.NRGBA))
	if err != nil {
		printAndExit("ERROR: could not decode image to NRGBA:", err)
	}

	decodedFilePath := strings.ReplaceAll(dataFilePath, ".encoded.png", "")
	if err = os.WriteFile(decodedFilePath, data, 0666); err != nil {
		printAndExit("ERROR: could not write decoded file:", err)
	}
}

func encodeData(dataFilePath string, imageSize int64) {
	data, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		printAndExit("ERROR: reading file"+dataFilePath+":", err)
	}

	imageCreator := image_creator.NewImageCreator(int(imageSize))
	imageCreator.GenImage()
	if err = imageCreator.Encode(data); err != nil {
		printAndExit("ERROR: encoding data in image:", err)
	}

	if err = imageCreator.SavePNG((dataFilePath) + ".encoded"); err != nil {
		printAndExit("ERROR: could not save image:", err)
	}
}

func printAndExit(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
}
