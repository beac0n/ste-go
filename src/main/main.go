package main

func main() {
	imgSize := 512

	imageCreator := NewImageCreator(imgSize)
	imageCreator.GetImage()
	imageCreator.SavePNG("image")
}
