package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"regexp"

	pikchr "github.com/jchenry/goldmark-pikchr"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	f, err := os.Open("README.md")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			&pikchr.Extender{},
		),
	)
	if err := md.Convert(b, &buf); err != nil {
		panic(err)
	}
	contents := buf.String()
	o, _ := os.Create("README.html")
	w := bufio.NewWriter(o)
	transformed := transformString(contents)
	w.WriteString(transformed)
	w.Flush()
	o.Close()
}

func transformString(input string) string {
	var re *regexp.Regexp
	var transformed string

	// Comment with text between 2 ticks
	// Example
	// <!-- `{% func () %}` -->
	re = regexp.MustCompile(`<!--\s*` + "`(.*?)`" + `\s*-->\n<(.*?)>`)
	transformed = re.ReplaceAllString(input, "$1\n<$2>")
	re = regexp.MustCompile(`<!--\s*` + "`(.*?)`" + `\s*-->`)
	transformed = re.ReplaceAllString(transformed, "$1")

	// Comment with text between 2 ticks
	// Example
	// <!-- class="bg-teal-700" -->
	re = regexp.MustCompile(`<!--(.*?)-->\n<(.*?)>`)
	transformed = re.ReplaceAllString(transformed, "<$2 $1>")
	re = regexp.MustCompile(`<!--(.*?)-->(.*?)<`)
	transformed = re.ReplaceAllString(transformed, "<div $1>$2<div><")
	return transformed
}
