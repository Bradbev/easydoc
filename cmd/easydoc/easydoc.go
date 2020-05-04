package main

import (
	"easydoc/internal/config"
	"easydoc/internal/markdown"
	"easydoc/internal/search"
	"easydoc/internal/walker"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	ignore "github.com/sabhiram/go-gitignore"
)

var root = flag.String("root", ".", "Set this to the root path to serve")
var port = flag.String("port", "localhost:8080", "Host and port to serve on")
var rootUrlFlag = flag.String("rootUrl", "http://localhost:8080", "This string will be prepended to make links fully qualified")

// rootUrl will take a copy of rootUrlFlag.  Exists only to remove the pointer deref noise
var rootUrl string

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
		liTags += T(`<li><a href onclick="tocClick('{{.base}}/#{{.file}}', '{{.base}}/{{.file}}')">{{.file}}</a></li>`, A{"file": file, "base": rootUrl})
	}
	return T(`
<!DOCTYPE html>
<html style="overflow: hidden;">

<head>
<title>Easydoc</title>
<link rel="stylesheet" type="text/css" href="{{.base}}/static/github-markdown.css">
<link rel="stylesheet" type="text/css" href="{{.base}}/static/easydoc.css">
<meta content="width=device-width, initial-scale=1" name="viewport">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="mobile-web-app-capable" content="yes">

<script
	src="https://code.jquery.com/jquery-1.12.4.min.js"
	integrity="sha256-ZosEbRLbNQzLpnKIkEdrPv7lOy9C27hHQ+Xp8a4MxAQ="
	crossorigin="anonymous">
</script>

<script>
function getBase() { return "{{.base}}" }
</script>

<script type="text/javascript" src="{{.base}}/static/easydoc.js"></script>
</head>

<body style="overflow: hidden;">
<header>
	<span id="application-title">Easydoc</span>
	<div id="search-field-container" role="search">
		<input id="search-field" type="search" placeholder="🔎 Search" required>
		<button id="search-field-button" type="button">🔎</button>
	</div>
</header>

<main>
	<nav id="toc">
		<ul>{{.li}}</ul>
		<div id="toc-drawer-handle">
			<button id="toc-drawer-toggle">
				<span class="icon-open">▶</span>
				<span class="icon-close">◀</span>
				<span class="icon-drawer-handle"></span>
				<span class="icon-close">◀</span>
				<span class="icon-open">▶</span>
			</button>
			<div id="toc-drawer-dragbar"></div>
		</div>
	</nav>

	<iframe id="content"></iframe>
</main>

</body>
</html>`,
		A{
			"li":   liTags,
			"base": rootUrl,
		})
}

func searchHtml(results []search.FileResult) string {
	hits := ""
	for _, r := range results {
		hits += T(`<h4><a class="searchlink" href="{{.base}}/{{.file}}">{{.file}}</a></h4><div class="results"><table>`, A{
			"file": r.File,
			"base": rootUrl,
		})
		for _, h := range r.Hits {
			hits += fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", h.LineNumber, html.EscapeString(h.Line))
		}
		hits += "</table></div>"
	}
	return T(`<html>
<head>
<link rel="stylesheet" type="text/css" href="{{.base}}/static/easydoc.css">

<script
	src="https://code.jquery.com/jquery-1.12.4.min.js"
	integrity="sha256-ZosEbRLbNQzLpnKIkEdrPv7lOy9C27hHQ+Xp8a4MxAQ="
	crossorigin="anonymous">
</script>

<script type="text/javascript" src="{{.base}}/static/easydoc.js"></script>
</head>

<body>
<div id="search">
<h1>Search Results</h1>
{{.hits}}
</div>
</body>
</html>`,
		A{
			"hits": hits,
			"base": rootUrl,
		})
}

const maxHitsPerPage = 5

func serveMarkdownFiles(searcher *search.Searcher, toc []string) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		// special case some of the paths
		fmt.Println("Handle", req.URL.Path)
		if req.URL.Path == "/" {
			values := req.URL.Query()
			toSearch := values.Get("search")
			if toSearch != "" {
				io.WriteString(w, searchHtml(searcher.Search(toSearch, maxHitsPerPage)))
				return
			}
		}
		if req.URL.Path == "/" || strings.ToLower(req.URL.Path) == "/index.html" {
			io.WriteString(w, tocHtml(toc))
			return
		}

		file := path.Join(*root, req.URL.String())

		result, err := markdown.MarkdownFileToHtml(file)
		if err != nil {
			result = "No file found at " + req.URL.String()
		}
		io.WriteString(w, result)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", handler)
	fmt.Println("Serving on ", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}

func makeIgnorer(conf *config.Config) (ignorer *ignore.GitIgnore) {
	ignoreFile := path.Join(*root, ".gitignore")
	fmt.Println(".gitignore at ", ignoreFile)
	var err error
	toIgnore := append(conf.Ignore, ".git")
	if _, err = os.Stat(ignoreFile); err == nil {
		fmt.Println("Found .gitignore")
		ignorer, err = ignore.CompileIgnoreFileAndLines(ignoreFile, toIgnore...)
	} else {
		ignorer, err = ignore.CompileIgnoreLines(toIgnore...)
	}
	if err != nil {
		log.Fatal(err)
	}
	return ignorer
}

func main() {
	flag.Parse()
	rootUrl = *rootUrlFlag
	conf := config.Load(path.Join(*root, "easydoc.json"))
	ignorerer := makeIgnorer(conf)
	markdown.SetUrlBase(*root, conf.ExternalUrlBase)

	searcher := search.Searcher{}
	fmt.Println("Scanning files at ", *root)
	files := walker.FindMarkdownFiles(ignorerer, *root)
	searcher.AddFiles(files)
	serveMarkdownFiles(&searcher, files)
}
