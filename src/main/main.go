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

	dataFilePath := flag.String("data-file-path", "", "path to file to encode in image(s) or path to image with data")
	encode := flag.Bool("encode", true, "encode data (default)")
	decode := flag.Bool("decode", false, "decode data")
	imageSize := flag.Int64("image-size", int64(512), "quarter size of image to encode data in")

	flag.Parse()

	if *dataFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*dataFilePath); os.IsNotExist(err) {
		fmt.Println("ERROR:", *dataFilePath, "does not exist")
		os.Exit(1)
	}

	if *encode && !*decode {
		encodeData(*dataFilePath, *imageSize)
	} else if *decode {
		decodeData(*dataFilePath)
	}
}

func decodeData(dataFilePath string) {
	f, _ := os.Open(dataFilePath)
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("ERROR: could not read file to image:", err)
		os.Exit(1)
	}
	data, err := image_creator.Decode(img.(*image.NRGBA))
	if err != nil {
		fmt.Println("ERROR: could not decode image to NRGBA:", err)
		os.Exit(1)
	}

	decodedFilePath := strings.ReplaceAll(dataFilePath, ".encoded.png", "")
	if err = os.WriteFile(decodedFilePath, data, 0666); err != nil {
		fmt.Println("ERROR: could not write decoded file:", err)
	}
}

func encodeData(dataFilePath string, imageSize int64) {
	data, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		fmt.Println("ERROR: reading file", dataFilePath, ":", err)
		os.Exit(1)
	}

	imageCreator := image_creator.NewImageCreator(int(imageSize))
	imageCreator.GenImage()
	if err = imageCreator.Encode(data); err != nil {
		fmt.Println("ERROR: encoding data in image:", err)
		os.Exit(1)
	}

	imageCreator.SavePNG((dataFilePath) + ".encoded")
}
