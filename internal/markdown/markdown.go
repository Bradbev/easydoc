package markdown

import (
	"io/ioutil"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/russross/blackfriday.v2"
)

func MarkdownFileToHtml(filename string) (string, error) {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return MarkdownStringToHtml(string(md)), nil
}

func MarkdownStringToHtml2(input string) string {
	html := blackfriday.Run([]byte(input), blackfriday.WithExtensions(blackfriday.CommonExtensions))
	return string(html)
}

func MarkdownStringToHtml(input string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.CompletePage
	opts := html.RendererOptions{
		Flags: htmlFlags,
		CSS:   "/static/github-markdown.css",
	}
	renderer := html.NewRenderer(opts)

	html := markdown.ToHTML([]byte(input), parser, renderer)
	return string(html)
}
