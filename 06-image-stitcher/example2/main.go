package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"math/rand"
	"net/http"
	"strconv"
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

// https://picsum.photos/id/%d/200
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	// first number is the ID, second number is the size

	n := r.URL.Query().Get("n")

	if n != "" {
		nInt, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}

		stitchHeight := 200 * nInt
		outputImage := image.NewRGBA(image.Rect(0, 0, 200, stitchHeight))

		for i := 0; i < nInt; i++ {
			fmt.Println(i)

			id := rand.Intn(250)

			url := "https://picsum.photos/id/" + strconv.Itoa(id) + "/200"
			fmt.Println(url)

			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close() // to avoid leaking memory / file descriptors

			img, _, err := image.Decode(resp.Body)
			if err != nil {
				panic(err)
			}

			// Draw image onto the output image
			draw.Draw(outputImage, image.Rect(0, 200*i, 200, 200*(i+1)), img, image.Point{}, draw.Over)
			//fmt.Fprintf(w, "Hello, %#v!", r.URL.Path[0:])
		}

		err = jpeg.Encode(w, outputImage, nil)
		if err != nil {
			fmt.Printf("Failed to encode output image: %v\n", err)
			return
		}

		//fmt.Fprintf(w, "body: %#v!", body)
		// 	if _, err := w.Write(body); err != nil {
		// 		panic(err)
		// 	}
	}
}
