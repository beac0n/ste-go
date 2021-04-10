package main

func main() {
	imgSize := 512

	imageCreator := NewImageCreator(imgSize)
	imageCreator.GenImage()
	imageCreator.SavePNG("image")
}
