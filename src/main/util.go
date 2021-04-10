package main

import "math/rand"

func GetRandomRGB() (uint8, uint8, uint8) {
	red := GetRandomColorValue()
	green := GetRandomColorValue()
	blue := GetRandomColorValue()
	return red, green, blue
}

func GetRandomColorValue() uint8 {
	return uint8(rand.Intn(256))
}
