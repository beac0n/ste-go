package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestGetImagePerformance(t *testing.T) {
	imageCreator := NewImageCreator(64)

	start := time.Now()
	imageCreator.GenImage()
	elapsed := time.Since(start)

	if elapsed > time.Millisecond*10 {
		t.Error("took too long:", elapsed)
	} else {
		log.Println("took", elapsed)
	}
}

func TestGetImageRandomness(t *testing.T) {
	testImgsFolder := "test_imgs"
	_ = os.RemoveAll(testImgsFolder)
	_ = os.Mkdir(testImgsFolder, 0700)

	imgSize := 1
	hashMap := make(map[string]struct{})

	collisions := 0

	rand.Seed(int64(1))

	for i := 0; i < imgSize*2*imgSize*2*255; i++ {
		imageCreator := NewImageCreator(imgSize)
		imageCreator.GenImage()
		hash := md5.New()
		hash.Write(imageCreator.img.Pix)

		sum := hash.Sum(nil)
		hexSum := hex.EncodeToString(sum)

		if _, isPresent := hashMap[hexSum]; isPresent {
			collisions += 1
			imageCreator.SavePNG(testImgsFolder + "/" + hexSum + "_" + strconv.Itoa(i))
			break
		} else {
			hashMap[hexSum] = struct{}{}
		}
	}

	if collisions > 0 {
		t.Error("got too many collisions:", collisions)
	} else {
		_ = os.RemoveAll(testImgsFolder)
	}

}
