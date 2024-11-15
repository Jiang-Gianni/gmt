package main

import (
	"flag"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Jiang-Gianni/gmt/markdown"
)

func main() {
	var out, dir, file, css, head string
	flag.StringVar(&out, "out", ".", "Output Directory")
	flag.StringVar(&dir, "dir", ".", "Input Directory")
	flag.StringVar(&file, "file", "", "File")
	flag.StringVar(&css, "css", "", "Css Link")
	flag.StringVar(&head, "head", "", "Header File")
	flag.Parse()

	var mdFiles []string
	if file != "" {
		mdFiles = append(mdFiles, file)
	} else {
		mdFiles = find(dir, ".md")
	}

	for _, mdFile := range mdFiles {
		htmlFile := out + "/" + strings.ReplaceAll(mdFile, ".md", ".html")
		c := markdown.Converter{
			MdFile:     mdFile,
			OutputFile: htmlFile,
			CssLink:    css,
			HeaderFile: head,
		}
		if err := c.ConvertFile(); err != nil {
			panic(err)
		}
	}

}

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}
