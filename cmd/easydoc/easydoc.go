package main

import (
	"easydoc/internal/config"
	"easydoc/internal/markdown"
	"easydoc/internal/walker"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	gitignore "github.com/sabhiram/go-gitignore"
)

var root = flag.String("root", ".", "Set this to the root path to serve")

// A is a templating helper
type A map[string]interface{}

// T is a templating helper
func T(input string, args A) string {
	w := &strings.Builder{}
	tmpl := template.Must(template.New("").Parse(input))
	err := tmpl.Execute(w, args)
	if err != nil {
		panic(err)
	}
	return w.String()
}

func tocHtml(toc []string) string {
	liTags := ""
	for _, file := range toc {
		liTags += T(`<li><a href="{{.file}}">{{.file}}</a></li>`, A{"file": file})
	}
	return T(`
<html>
<head>
<link rel="stylesheet" type="text/css" href="/static/github-markdown.css">
<link rel="stylesheet" type="text/css" href="/static/easydoc.css">

<script
  src="https://code.jquery.com/jquery-1.12.4.min.js"
  integrity="sha256-ZosEbRLbNQzLpnKIkEdrPv7lOy9C27hHQ+Xp8a4MxAQ="
  crossorigin="anonymous">
</script>

<script type="text/javascript" src="static/easydoc.js"></script>
</head>
<body>

<div id="toc">
	<ul>{{.li}}</ul>
</div>

<iframe id="content"></iframe>

</body>
</html>`,
		A{"li": liTags})
}

func serveMarkdownFiles(toc []string) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		// ugh, special case some of the paths
		if req.URL.Path == "/" || strings.ToLower(req.URL.Path) == "/index.html" {
			io.WriteString(w, tocHtml(toc))
			return
		}

		file := path.Join(*root, req.URL.String())

		result, err := markdown.MarkdownFileToHtml(file)
		if err != nil {
			result = markdown.MarkdownStringToHtml("", "### No file found at "+req.URL.String())
		}
		io.WriteString(w, result)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", handler)
	fmt.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func makeIgnorer(conf *config.Config) (ignorer *gitignore.GitIgnore) {
	ignoreFile := path.Join(*root, ".gitignore")
	var err error
	if _, err = os.Stat(ignoreFile); err == nil {
		ignorer, err = gitignore.CompileIgnoreFileAndLines(ignoreFile, conf.Ignore...)
	}
	ignorer, err = gitignore.CompileIgnoreLines(conf.Ignore...)
	if err != nil {
		log.Fatal(err)
	}
	return ignorer
}

func main() {
	flag.Parse()
	conf := config.Load(path.Join(*root, "easydoc.json"))
	ignorerer := makeIgnorer(conf)
	markdown.SetUrlBase(conf.ExternalUrlBase)

	fmt.Println("Scanning files at ", *root)
	files := walker.FindMarkdownFiles(ignorerer, *root)
	serveMarkdownFiles(files)
}
