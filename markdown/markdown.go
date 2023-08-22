package markdown

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Jiang-Gianni/gmt/css"
	pikchr "github.com/jchenry/goldmark-pikchr"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var Md = goldmark.New(
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

func ConvertFile(mdFile string, outputFile string, css string) {
	f, err := os.Open(mdFile)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	var buf bytes.Buffer
	if err := Md.Convert(b, &buf); err != nil {
		panic(err)
	}
	contents := buf.String()
	os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	o, _ := os.Create(outputFile)
	w := bufio.NewWriter(o)
	transformed := transformString(contents, css)
	w.WriteString(transformed)
	w.Flush()
	o.Close()
}

func transformString(input string, cssLink string) string {
	var re *regexp.Regexp
	var transformed string

	// Comment with text between 2 ticks
	// Example
	// <!-- `{% func () %}` -->
	re = regexp.MustCompile(`<!--\s*` + "`(.*?)`" + `\s*-->\n<(.*?)>`)
	transformed = re.ReplaceAllString(input, "$1\n<$2>")
	re = regexp.MustCompile(`<!--\s*` + "`(.*?)`" + `\s*-->`)
	transformed = re.ReplaceAllString(transformed, "$1")

	// Comment with attributes to be injected
	// Example
	// <!-- class="bg-teal-700" -->
	re = regexp.MustCompile(`<!--(.*?)-->\n<(.*?)>`)
	transformed = re.ReplaceAllString(transformed, "<$2 $1>")
	re = regexp.MustCompile(`<!--(.*?)--><(.*?)>`)
	transformed = re.ReplaceAllString(transformed, "<$2 $1>")
	re = regexp.MustCompile(`><!--(.*?)-->(.*?)</`)
	transformed = re.ReplaceAllString(transformed, "><div $1>$2</div><")

	// If the element alreadi had the same class attribute
	// Example (after the previous injection)
	// <div class="pikchr-svg" class="bg-teal-700">
	re = regexp.MustCompile(`<(.*?) class=\"(.*?)\"(.*?)class=\"(.*?)\"(.*?)>`)
	transformed = re.ReplaceAllString(transformed, `<$1 class="$2 $4" $3 $5>`)

	classes := css.GetClasses(transformed)
	styles := css.GetStyles(classes)
	transformed = addStyleTag(styles, cssLink, transformed)
	return transformed
}
