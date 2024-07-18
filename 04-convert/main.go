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

	lines := strings.Split(s, "\n")
	keys := strings.Split(lines[0], ",")
	hms := []map[string]string{}
	for _, valueLine := range lines[1:] {
		if valueLine == "" {
			break
		}
		values := strings.Split(valueLine, ",")
		hm := map[string]string{}
		for j, key := range keys {
			hm[key] = values[j]
		}
		hms = append(hms, hm)
	}

	if isShow {
		fmt.Println("--- Showing", filename)
		for _, hm := range hms {
			fmt.Printf("%#v\n", hm)
		}
	}
	if isCount {
		fmt.Printf("--- Number of rows: %d\n", len(hms))
	}
}

func readFileToString(filename string) string {
	bs, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
