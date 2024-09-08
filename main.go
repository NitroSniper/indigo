package main

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"os"
)

func e(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	source, err := os.ReadFile("./example.md")
	e(err)

	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}

	fmt.Print(string(buf.Bytes()))
}
