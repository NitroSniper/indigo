package main

import (
	// "bytes"
	"fmt"
	// "github.com/yuin/goldmark"
	"os"
)

func e(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Hello World")
	bytes, err := os.ReadFile("./example.md")
	e(err)
	fmt.Print(string(bytes))

	// var buf bytes.Buffer
	// if err := goldmark.Convert(source, &buf); err != nil {
	// 	panic(err)
	// }
}
