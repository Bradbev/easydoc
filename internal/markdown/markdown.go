package markdown

import (
	"io/ioutil"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/russross/blackfriday.v2"
)

var urlBase string

func SetUrlBase(base string) {
	urlBase = base
}

func MarkdownFileToHtml(filename string) (string, error) {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return MarkdownStringToHtml(filename, string(md)), nil
}

func MarkdownStringToHtml2(input string) string {
	html := blackfriday.Run([]byte(input), blackfriday.WithExtensions(blackfriday.CommonExtensions))
	return string(html)
}

func MarkdownStringToHtml(filename, input string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)

	raw := markdown.ToHTML([]byte(input), parser, renderer)
	html := `
	<html> 
	<link rel="stylesheet" type="text/css" href="/static/github-markdown.css">
	<link rel="stylesheet" type="text/css" href="/static/easydoc.css">
	<body>
	`
	if urlBase != "" {
		html += `<div class="header">View file externally at <a target="_blank" href="` + urlBase + filename + `">` + filename + `</a><HR></div>`
	}
	html += `<div class="markdown-body">` + string(raw) + `</div></body> </html>`
	return html
}
