package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

func FindMarkdownFiles(ignorer *gitignore.GitIgnore, base string) []string {
	result := make([]string, 0)
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if ignorer.MatchesPath(path) && info.IsDir() {
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

	return result
}
