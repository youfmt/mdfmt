package mdfmt

import (
	"strings"
	"testing"

	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
)

func (b block) Equals(other interface{}) bool {
	switch other := other.(type) {
	case block:
		if b.blockType != other.blockType {
			return false
		}

		for i, _ := range b.lines {
			if b.lines[i] != other.lines[i] {
				return false
			}
		}

		return true

	default:
	}

	return false
}

func describeBlockLexer(c gospec.Context) {
	c.Specify("A blockLexer", func() {
		c.Specify("Emits a paragraph", func() {
			blocks := newBlockLexer(strings.NewReader(`
some paragraph
text
`)).Run()

			b := <-blocks
			c.Expect(b, Equals, block{
				blockType: BT_PARAGRAPH,
				lines: []string{
					"some paragraph",
					"text",
				},
			})

			b = <-blocks
			c.Expect(b, Equals, block{
				blockType: BT_ERROR,
				lines: []string{
					"EOF",
				},
			})
		})
	})
}

func TestBlockLexer(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(describeBlockLexer)
	gospec.MainGoTest(r, t)
}
