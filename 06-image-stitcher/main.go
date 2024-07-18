package main

import (
	"fmt"
	"io"
	"net/http"
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
	resp, err := http.Get("https://picsum.photos/id/237/600")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // to avoid leaking memory / file descriptors
	body, err := io.ReadAll(resp.Body)
	//fmt.Fprintf(w, "body: %#v!", body)
	if _, err := w.Write(body); err != nil {
		panic(err)
	}
	//fmt.Fprintf(w, "Hello, %#v!", r.URL.Path[0:])
}
