package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	s := Searcher{
		files:  []string{"a", "b"},
		sorted: true,
		content: map[string][]string{
			"a": []string{"a", "a", "c"},
			"b": []string{"b", "b", "c"},
		},
	}

	a := s.Search("a", 0)
	assert.Len(t, a, 1)
	assert.Equal(t, "a", a[0].File)
	assert.Len(t, a[0].Hits, 2)

	ab := s.Search("b|a", 0)
	assert.Len(t, ab, 2)
}
