package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func hello(w http.ResponseWriter, req *http.Request) {
	buf := foo()

	templ := template.Must(template.New("index").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Exmpale</title>
	</head>
	<body>
		{{.Markdown}}
	</body>
	</html>
`))
	data := struct {
		Markdown template.HTML
	}{
		Markdown: template.HTML(buf.Bytes()),
	}

	if err := templ.Execute(w, data); err != nil {
		panic(err)
	}
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func foo() bytes.Buffer {
	source, err := os.ReadFile("./example.md")
	if err != nil {
		panic(err)
	}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
	return buf
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/headers", headers)
	log.Fatal(http.ListenAndServe(":8090", mux))
}
