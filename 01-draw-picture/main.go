package main

import (
	"fmt"
	"strconv"
)

// string to int -> can error
// int to string -> won't error

func main() {
	// converting a number into a string, and printing it
	fmt.Printf("%#v\n", strconv.Itoa(42))

	// converting a string into a number
	//convertedInt, err := strconv.Atoi("4asd2asdf")
	//if err == nil {
	if convertedInt, err := strconv.Atoi("4asd2asdf"); err == nil {
		fmt.Printf("%#v\n", convertedInt)
	} else {
		panic(err)
		//fmt.Printf("unable to convert int: %s\n", err.Error())
	}
}
