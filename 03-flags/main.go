package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var isShow bool
	var isCount bool
	flag.BoolVar(&isCount, "count", false, "count lines mode")
	flag.BoolVar(&isShow, "show", false, "show file mode")
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		fmt.Println("must supply a filename")
		os.Exit(1)
	}
	s := readFileToString(filename)

	if isShow {
		fmt.Println("--- Showing", filename)
		fmt.Println(s)
	}
	if isCount {
		lines := strings.Split(s, "\n")
		fmt.Printf("--- Number of lines: %d\n", len(lines))
	}
}

func readFileToString(filename string) string {
	bs, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
