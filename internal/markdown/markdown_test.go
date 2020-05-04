package markdown

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestFixUrls(t *testing.T) {
	raw := []byte(`start [link](/addr) 
second [link2](../addr2)
`)

	fixed := fixInternalLinks("base:/", "nested/n2/file.md", raw)

	assert.Equal(t, `start [link](base://addr) 
second [link2](base://nested/addr2)
`, string(fixed))
}
