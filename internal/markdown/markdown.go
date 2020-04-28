package markdown

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var urlBase string
var rootPath string

func SetUrlBase(rootPathArg, base string) {
	urlBase = base
	rootPath = rootPathArg
}

func MarkdownFileToHtml(filename string) (string, error) {
	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	gold := goldmark.New(goldmark.WithExtensions(extension.GFM))
	var buf bytes.Buffer
	if err = gold.Convert(filedata, &buf); err != nil {
    	return "", err
	}

	html := `
	<html> 
	<link rel="stylesheet" type="text/css" href="/static/github-markdown.css">
	<link rel="stylesheet" type="text/css" href="/static/easydoc.css">
	<body>
	`
	if urlBase != "" {
		if strings.HasPrefix(filename, rootPath) {
			filename = filename[len(rootPath):]
		}
		html += `<div class="header">View file externally at <a target="_blank" href="` + urlBase + filename + `">` + filename + `</a><HR></div>`
	}
	html += `<div class="markdown-body">` + buf.String() + `</div></body> </html>`
	return html, nil
}