package markdown

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Jiang-Gianni/gmt/css"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	pikchr "github.com/jchenry/goldmark-pikchr"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"

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
		highlighting.NewHighlighting(
			highlighting.WithStyle("hrdark"),
			highlighting.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
			),
		),
	),
)

type Converter struct {
	MdFile     string
	OutputFile string
	CssLink    string
	HeaderFile string
}

func (c *Converter) ConvertFile() error {
	mdBytes, err := os.ReadFile(c.MdFile)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := Md.Convert(mdBytes, &buf); err != nil {
		return err
	}

	contents := buf.String()
	if err := os.MkdirAll(filepath.Dir(c.OutputFile), os.ModePerm); err != nil {
		return err
	}

	o, err := os.Create(c.OutputFile)
	if err != nil {
		return err
	}
	defer o.Close()

	if c.HeaderFile != "" {
		headBytes, err := os.ReadFile(c.HeaderFile)
		if err != nil {
			return err
		}
		_, err = o.Write(headBytes)
		if err != nil {
			return err
		}
	}

	w := bufio.NewWriter(o)
	transformed := transformString(contents, c.CssLink)
	_, err = w.WriteString(transformed)
	if err != nil {
		return err
	}
	return w.Flush()
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
	re = regexp.MustCompile(`(\w+)<!--(.*?)-->(\w+)`)
	transformed = re.ReplaceAllString(transformed, "<p $2>$1</p>$3")
	re = regexp.MustCompile(`<!--(.*?)-->\n<ul>`)
	transformed = re.ReplaceAllString(transformed, "<ul $1>")
	re = regexp.MustCompile(`<!--(.*?)class=\"(.*?)\"(.*?)-->\r?\n<div(.*?) class=\"(.*?)\"(.*?)>`)
	transformed = re.ReplaceAllString(transformed, `<div$4 class="$5 $2" $1 $3 $6>`)
	re = regexp.MustCompile(`<([^<>]*)>([^<>]*)</(.*?)><!--(.*?)-->`)
	transformed = re.ReplaceAllString(transformed, "<$1 $4>$2</$3>")
	re = regexp.MustCompile(`<p>((^[<>/]|.)*?)</p>\r?\n<!--(.*?)-->`)
	transformed = re.ReplaceAllString(transformed, "<p $3>$1</p>")
	re = regexp.MustCompile(`<(.*?)>(.*?)</(.*?)>\r?\n<!--(.*?)-->`)
	transformed = re.ReplaceAllString(transformed, "<$1 $4>$2</$1>")

	classes := css.GetClasses(transformed)
	styles := css.GetStyles(classes)
	transformed = addStyleTag(styles, cssLink, transformed)
	return transformed
}
