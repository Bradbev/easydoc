package walker

import (
	"fmt"
	"testing"

	gitignore "github.com/sabhiram/go-gitignore"
)

func TestGitIgnore(t *testing.T) {
	ignoreObject, error := gitignore.CompileIgnoreLines([]string{"node_modules", "*.out", "foo/*.c"}...)
	if error != nil {
		panic("Error when compiling ignore lines: " + error.Error())
	}

	// You can test the ignoreObject against various paths using the
	// "MatchesPath()" interface method. This pretty much is up to
	// the users interpretation. In the case of a ".gitignore" file,
	// a "match" would indicate that a given path would be ignored.
	fmt.Println(ignoreObject.MatchesPath("node_modules/test/foo.js"))
	fmt.Println(ignoreObject.MatchesPath("node_modules2/test.out"))
	fmt.Println(ignoreObject.MatchesPath("test/foo.js"))
	t.Fail()
}
