package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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

func getImagesSerial(imageNum int) []image.Image {
	images := []image.Image{}
	for i := 0; i < imageNum; i++ {
		images = append(images, getImage())
	}
	return images
}

// each goroutine will fetch an image individually (getImage())
// 		- put that into results channel
//		- all will happen in the background

// read off the results channel

func getImagesParallel(imageNum int) []image.Image {
	// JS (single threaded)
	// await Promise.all([p1, p2, p3])

	// goroutines
	// channels

	images := []image.Image{}
	results := make(chan image.Image)
	for i := 0; i < imageNum; i++ {
		fn := func() {
			results <- getImage()
		}
		go fn()
	}

	for i := 0; i < imageNum; i++ {
		image := <-results // reading image from the channel
		images = append(images, image)
	}

	return images
}

func getImage() image.Image {
	rand.Seed(time.Now().UnixNano())

	min := 1
	max := 600

	randomNum := rand.Intn(max-min+1) + min

	url := fmt.Sprintf("https://picsum.photos/id/%d/600", randomNum)

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	img, _, err := image.Decode(bytes.NewReader(body))

	if err != nil {
		//panic(err)
		return getImage()
	}

	return img
}

func getBytes(image image.Image) []byte {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func stitchImage(image1 image.Image, image2 image.Image, vertical bool, start int) image.Image {
	// https://stackoverflow.com/questions/35964656/golang-how-to-concatenate-append-images-to-one-another

	var rgba *image.RGBA

	if !vertical {
		// horizontal
		sp2 := image.Point{0 + start, 0}
		r2 := image.Rectangle{sp2, sp2.Add(image2.Bounds().Size())}
		r := image.Rectangle{image.Point{0, 0}, image.Point{2400, r2.Max.Y}}

		rgba = image.NewRGBA(r)
		draw.Draw(rgba, image1.Bounds(), image1, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, r2, image2, image.Point{0, 0}, draw.Src)
	} else {
		// vertical
		sp2 := image.Point{0, image1.Bounds().Dy()}
		r2 := image.Rectangle{sp2, sp2.Add(image2.Bounds().Size())}
		r := image.Rectangle{image.Point{0, 0}, r2.Max}

		rgba = image.NewRGBA(r)
		draw.Draw(rgba, image1.Bounds(), image1, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, r2, image2, image.Point{0, 0}, draw.Src)
	}

	return rgba
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	st := time.Now()

	// Extract the image number from the URL path
	path := r.URL.Path
	if path == "/favicon.ico" {
		return
	}
	path = strings.TrimPrefix(path, "/")
	imageNum, err := strconv.Atoi(path)
	if err != nil {
		imageNum = 1
	}

	images := getImagesParallel(imageNum)
	//images := getImagesSerial(imageNum)

	var vertical image.Image
	var horizontal image.Image

	vertical = image.NewRGBA(image.Rect(0, 0, 0, 0))
	horizontal = image.NewRGBA(image.Rect(0, 0, 0, 0))

	for i, v := range images {
		if i%4 == 0 && i != 0 {
			vertical = stitchImage(vertical, horizontal, true, (i%4)*600)
			horizontal = image.NewRGBA(image.Rect(0, 0, 0, 0))
			horizontal = stitchImage(horizontal, v, false, (i%4)*600)
		} else {
			horizontal = stitchImage(horizontal, v, false, (i%4)*600)
		}
	}

	vertical = stitchImage(vertical, horizontal, true, 0)

	if _, err := w.Write(getBytes(vertical)); err != nil {
		panic(err)
	}

	fmt.Printf("--- Request took %.2f seconds\n", time.Since(st).Seconds())
}
