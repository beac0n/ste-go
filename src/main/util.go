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

func ReverseValues(a, b int) (int, int) {
	oldA := a
	a = b
	b = oldA
	return a, b
}
