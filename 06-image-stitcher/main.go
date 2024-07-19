package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

/// Assignment:
/// Change ImageHandler to accept:
///  /N - where N is a number
///
/// - Grab N images
/// - Stitch them together using https://pkg.go.dev/image
///
/// Stretch:
/// - Do the fetches in parallel

func main() {
	http.HandleFunc("/", ImageHandler)
	fmt.Println("Listening on port http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// benchmarker programs
// Apache Bench -- `ab`
// Python Swarm
// "dark launch" -- feature flags
//		- make request to the old version
//			- the real data
//		- make a request to the new version
//          - a percentage

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

// range 0 - 1029 (inclusive)
func grabImage() image.Image {
	imageID := randRange(0, 1029+1)
	url := fmt.Sprintf("https://picsum.photos/id/%d/200", imageID)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		// some of the images we can't decode for some reason, so... try again in a janky way
		return grabImage()
	}
	return img
}

func grabNImages(n int) []image.Image {
	res := make([]image.Image, n)

	results := make(chan image.Image)
	for i := 0; i < n; i++ {
		fn := func() {
			image := grabImage()
			results <- image
		}
		go fn()
	}
	for i := 0; i < n; i++ {
		res[i] = <-results
	}
	return res
}

// https://picsum.photos/id/%d/200
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	st := time.Now()

	numImages, err := strconv.Atoi(r.URL.Path[1:])
	if err != nil {
		fmt.Fprintf(w, "Path must be a number")
		return
	}

	// serial
	//images := make([]image.Image, numImages)
	//for i, _ := range images {
	//	images[i] = grabImage()
	//}

	// parallel
	images := grabNImages(numImages)

	finalWidth := 200
	finalHeight := 200 * len(images)
	finalImg := image.NewRGBA(image.Rect(0, 0, finalWidth, finalHeight))
	for i, img := range images {
		x := 0
		y := i * 200
		draw.Draw(finalImg, image.Rect(0, i*200, x+200, y+200), img, image.Point{0, 0}, draw.Over)
	}

	if err := png.Encode(w, finalImg); err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	fmt.Printf("Got %d images in %.2f seconds\n", numImages, time.Since(st).Seconds())
}
