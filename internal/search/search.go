package search

import (
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

type Searcher struct {
	files   []string
	content map[string][]string
	sorted  bool
}

func (s *Searcher) AddFile(file string) {
	s.sorted = false
	s.files = append(s.files, file)
}

func (s *Searcher) AddFiles(files []string) {
	s.sorted = false
	s.files = append(s.files, files...)
}

func (s *Searcher) Search(regex string, maxResults int) []FileResult {
	if !s.sorted {
		s.sorted = true
		sort.Sort(sort.StringSlice(s.files))
	}
	return s.getResults(regex, maxResults)
}

func (s *Searcher) getResults(regex string, maxResults int) []FileResult {
	if s.content == nil {
		s.content = make(map[string][]string)
	}
	var result []FileResult
	r, _ := regexp.Compile("(?i)" + regex)
	for _, file := range s.files {
		lines, ok := s.content[file]
		if !ok {
			content, _ := ioutil.ReadFile(file)
			if content != nil {
				lines = strings.Split(string(content), "\n")
				s.content[file] = lines
			}
		}
		thisResult := FileResult{}
		for i, line := range lines {
			if r.FindString(line) != "" {
				thisResult.File = file
				thisResult.Hits = append(thisResult.Hits, Hit{i + 1, line})
				if maxResults > 0 && len(thisResult.Hits) > maxResults {
					thisResult.Hits = append(thisResult.Hits, Hit{i + 2, "...more hits..."})
					break
				}
			}
		}
		if thisResult.File != "" {
			result = append(result, thisResult)
		}
	}
	return result
}

type Hit struct {
	LineNumber int
	Line       string
}

type FileResult struct {
	File string
	Hits []Hit
}
