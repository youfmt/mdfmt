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

	assumeEOFEmitted := func(blocks <-chan block) {
		b := <-blocks
		c.Assume(b, Equals, block{
			blockType: BT_ERROR,
			lines: []string{
				"EOF",
			},
		})

	}
	c.Specify("A blockLexer", func() {

		c.Specify("Emits a paragraph", func() {
			blocks := newBlockLexer(strings.NewReader("some paragraph\ntext")).Run()
			defer assumeEOFEmitted(blocks)

			b := <-blocks
			c.Expect(b, Equals, block{
				blockType: BT_PARAGRAPH,
				lines: []string{
					"some paragraph",
					"text",
				},
			})

		})

		c.Specify("Emits several paragraphs", func() {
			blocks := newBlockLexer(strings.NewReader("paragraph1\n\nparagraph2")).Run()
			defer assumeEOFEmitted(blocks)

			b := <-blocks
			c.Expect(b, Equals, block{
				BT_PARAGRAPH,
				[]string{"paragraph1"},
			})

			b = <-blocks
			c.Expect(b, Equals, block{
				BT_PARAGRAPH,
				[]string{"paragraph2"},
			})
		})
	})
}

func TestBlockLexer(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(describeBlockLexer)
	gospec.MainGoTest(r, t)
}
