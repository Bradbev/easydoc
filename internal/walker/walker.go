package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/facebookgo/symwalk"
	ignore "github.com/sabhiram/go-gitignore"
)

func FindMarkdownFiles(ignorer *ignore.GitIgnore, base string) []string {
	return FindFiles(ignorer, base, `.*\.md$`)
}
func FindFiles(ignorer *ignore.GitIgnore, base string, regex string) []string {
	matcher := regexp.MustCompile(regex)
	result := make([]string, 0)
	startTime := time.Now()
	err := symwalk.Walk(base, func(path string, info os.FileInfo, err error) error {
		path = strings.ReplaceAll(path, "\\", "/")
		relativePath := strings.TrimPrefix(path, base)
		// fmt.Println(relativePath)
		if ignorer.MatchesPath(relativePath) && info.IsDir() {
			fmt.Println("Skipping dir:", relativePath)
			return filepath.SkipDir
		}

		if err != nil {
			fmt.Printf("failed to access path %q: %v\n", path, err)
			return nil
		}
		if matcher.MatchString(strings.ToLower(path)) {
			stripped := strings.TrimPrefix(path, base)
			result = append(result, stripped)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", base, err)
		panic("")
	}
	fmt.Println("Scan took ", time.Now().Sub(startTime))

	return result
}
