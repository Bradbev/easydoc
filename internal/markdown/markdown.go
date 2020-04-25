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

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)

	raw := markdown.ToHTML([]byte(input), parser, renderer)
	html := `
<html>
	<link rel="stylesheet" type="text/css" href="/static/github-markdown.css">
	<body class="markdown-body">` + string(raw) + `
	</body>
</html>`
	return html
}
