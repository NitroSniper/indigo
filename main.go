package main

import (
	"bytes"
	_ "embed" // embed is not used directly but for it's macro
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

//go:embed assets/template/base.html
var base_template string

//go:embed assets/flavor/github.css
var flavor string

func hello(w http.ResponseWriter, req *http.Request) {
	markdown := foo()
	boxing := `
	.markdown-body {
		padding: 64px;
	}
	body {
	  margin: 0;
	}
	`
	templ := template.Must(template.New("index").Parse(base_template))
	data := struct {
		Markdown template.HTML
		Flavor   template.CSS
	}{
		Markdown: template.HTML(markdown.Bytes()),
		Flavor:   template.CSS(flavor + boxing),
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
