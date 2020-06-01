package markdown

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
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
	if strings.HasPrefix(filename, hostFileBase) {
		filename = filename[len(hostFileBase):]
	}
	filedata = fixInternalLinks(hostUrlBase, filename, filedata)

	gold := goldmark.New(goldmark.WithExtensions(extension.GFM, highlighting.Highlighting))
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
		html += `<div class="header">View file externally at <a target="_blank" href="` + externalUrl + filename + `">` + filename + `</a><HR></div>`
	}
	html += `<div class="markdown-body">` + buf.String() + `</div></body> </html>`
	return html, nil
}

var linkRegex = regexp.MustCompile(`\[(?P<name>.*?)\]\s*\((?P<link>.*?)\)`)

//                   match against [name](link), with any space between ]()

func fixInternalLinks(hostUrlBase, currentFile string, data []byte) []byte {
	return linkRegex.ReplaceAllFunc(data, func(data []byte) []byte {
		matches := linkRegex.FindStringSubmatch(string(data))
		name, link := matches[1], matches[2]
		if strings.HasPrefix(link, "http") {
			return data
		}
		if strings.HasPrefix(link, "/") {
			fullLink := hostUrlBase + link
			return []byte(fmt.Sprintf("[%s](%s)", name, fullLink))
		}
		fileBase := path.Dir(currentFile)
		relLink := hostUrlBase + "/" + path.Clean(path.Join(fileBase, link))
		return []byte(fmt.Sprintf("[%s](%s)", name, relLink))
	})
}
