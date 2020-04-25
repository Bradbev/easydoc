package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindMarkdownFiles(base string) []string {
	result := make([]string, 0)
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			panic("")
		}
		if strings.HasSuffix(strings.ToLower(path), ".md") {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", base, err)
		panic("")
	}

	return result
}
