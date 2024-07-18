package main

import "os"
import "fmt"
import "strings"

func main() {
  s := readFileToString("some.csv")
  lines := strings.Split(s, "\n")
  headerKeys := strings.Split(lines[0], ",")
  values := strings.Split(lines[1], ",")
  fmt.Printf("headerKeys: %#v\n", headerKeys)
  fmt.Printf("values: %#v\n", values)
  hm := map[string]string{}
  for i, key := range headerKeys {
    hm[key] = values[i]
  }
  fmt.Printf("collected into hm: %#v\n", hm)
}

func readFileToString(filename string) string {
	bs, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

