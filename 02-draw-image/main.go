package main

import (
	"bufio"
	"image"
	"image/png"
	"os"
)

type Image struct{}

func main() {
	m := Image{}
	saveImage(m, "image.png")
}

func saveImage(img image.Image, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	err = png.Encode(w, img)
	if err != nil {
		panic(err)
	}
}
