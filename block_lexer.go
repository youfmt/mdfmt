//go:generate stringer -type=blockType
package mdfmt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type blockType int

const (
	BT_ERROR blockType = iota

	// Containter blocks
	BT_BLOCK_QUOTE

	BT_LIST_ITEM
	BT_LIST

	// Leaf blocks
	BT_ATX_HEADER
	BT_SETEXT_HEADER

	BT_CODE_BLOCK

	BT_PARAGRAPH
)

type lexerStateFn func(*blockLexer) lexerStateFn

type blockLexer struct {
	input *bufio.Reader

	lines []string

	blocks chan<- block
}

func newBlockLexer(input io.Reader) *blockLexer {
	buf := bufio.NewReader(input)

	return &blockLexer{input: buf}
}

type block struct {
	blockType blockType
	lines     []string
}

func (b block) String() string {
	return fmt.Sprintf("{%v, [%s]}", b.blockType, strings.Join(b.lines, ", "))
}

// run emits blocks
func (bl *blockLexer) Run() <-chan block {
	blocks := make(chan block)
	go func() {
		for stateFn := lexBlock; stateFn != nil; {
			stateFn = stateFn(bl)
		}
		close(blocks)
	}()

	bl.blocks = blocks
	return blocks
}

func (bl *blockLexer) peek() (rune, error) {
	r, _, err := bl.input.ReadRune()
	bl.input.UnreadRune()
	return r, err
}

func (bl *blockLexer) emit(bt blockType) {
	lines := bl.lines
	bl.lines = nil

	bl.blocks <- block{
		bt,
		lines,
	}
}

func lexBlock(l *blockLexer) lexerStateFn {
	_, err := l.peek()
	if err != nil {
		l.lines = append(l.lines, err.Error())
		l.emit(BT_ERROR)
		return nil
	}

	return lexParagraph
}

func lexBlockQuote(*blockLexer) lexerStateFn {
	return nil
}

func lexListItem(*blockLexer) lexerStateFn {
	return nil
}

func lexList(*blockLexer) lexerStateFn {
	return nil
}

func lexHeader(*blockLexer) lexerStateFn {
	return nil
}

func lexCodeBlock(*blockLexer) lexerStateFn {
	return nil
}

func lexParagraph(l *blockLexer) lexerStateFn {
	var line []byte
	var err error

	for err == nil {
		line, _, err = l.input.ReadLine()
		if string(line) == "" {
			continue
		}
		l.lines = append(l.lines, string(line))
	}

	l.emit(BT_PARAGRAPH)

	return lexBlock
}
