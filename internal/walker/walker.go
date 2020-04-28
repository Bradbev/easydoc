package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ignore "github.com/sabhiram/go-gitignore"
)

func FindMarkdownFiles(ignorer *ignore.GitIgnore, base string) []string {
	result := make([]string, 0)
	startTime := time.Now()
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		path = strings.ReplaceAll(path, "\\", "/")
		relativePath := strings.TrimPrefix(path, base)
		// fmt.Println(relativePath)
		if ignorer.MatchesPath(relativePath) && info.IsDir() {
			fmt.Println("Skipping dir:", relativePath)
			return filepath.SkipDir
		}

		if err != nil {
			fmt.Printf("failured to access path %q: %v\n", path, err)
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".md") {
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
