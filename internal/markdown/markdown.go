package markdown

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var externalUrl string
var hostUrlBase string
var hostFileBase string

func SetUrlBase(hostFileBaseArg, hostUrl, externalUrlArg string) {
	externalUrl = externalUrlArg
	hostUrlBase = hostUrl
	hostFileBase = hostFileBaseArg
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
	<link rel="stylesheet" type="text/css" href="` + hostUrlBase + `/static/github-markdown.css">
	<link rel="stylesheet" type="text/css" href="` + hostUrlBase + `/static/easydoc.css">
	<body>
	`
	if externalUrl != "" {
		if strings.HasPrefix(filename, hostFileBase) {
			filename = filename[len(hostFileBase):]
		}
		html += `<div class="header">View file externally at <a target="_blank" href="` + externalUrl + filename + `">` + filename + `</a><HR></div>`
	}
	html += `<div class="markdown-body">` + buf.String() + `</div></body> </html>`
	return html, nil
}
