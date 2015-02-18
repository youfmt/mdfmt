//go:generate stringer -type=blockType
package mdfmt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
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
	var lines string
	if len(b.lines) == 0 {
		lines = "EMPTY"
	} else {
		lines = strings.Join(b.lines, ", ")
	}

	return fmt.Sprintf("{%v, [%s]}", b.blockType, lines)
}

// run emits blocks
func (l *blockLexer) Run() <-chan block {
	blocks := make(chan block)
	go func() {
		for stateFn := lexBlock; stateFn != nil; {
			stateFn = stateFn(l)
		}
		close(blocks)
	}()

	l.blocks = blocks
	return blocks
}

func (l *blockLexer) peek() (rune, error) {
	r, _, err := l.input.ReadRune()
	if err != nil {
		return r, err
	}
	return r, l.input.UnreadRune()
}

func (l *blockLexer) emit(bt blockType) {
	lines := l.lines
	l.lines = nil

	l.blocks <- block{
		bt,
		lines,
	}
}

func (l *blockLexer) emitError(err error) {
	l.lines = append(l.lines, err.Error())
	l.emit(BT_ERROR)
}

func (l *blockLexer) consumeLine() error {
	line, _, err := l.input.ReadLine()
	if err != nil {
		return err
	}
	l.lines = append(l.lines, string(line))
	return nil
}

func (l *blockLexer) immolateLine() error {
	_, _, err := l.input.ReadLine()
	return err
}

func (l *blockLexer) immolateRune() error {
	_, _, err := l.input.ReadRune()
	return err
}

var ErrUnexpectedInput = errors.New("unexpected input")

func lexBlock(l *blockLexer) lexerStateFn {
	r, err := l.peek()
	if err != nil {
		l.emitError(err)
		return nil
	}

	switch {
	case isWhiteSpace(r):
		return startLexingSpace(l)

	case isParagraph(r):
		return lexParagraph
	case isBlockQuote(r):
		return lexBlockQuote

	default:
		l.emitError(ErrUnexpectedInput)
		return nil
	}
}

func (l *blockLexer) peekLine() ([]byte, error) {
	i := bufio.NewReader(l.input)
	line, _, err := i.ReadLine()
	return line, err
}

func startLexingSpace(l *blockLexer) lexerStateFn {
	var fn lexerStateFn

	fn = func(l *blockLexer) lexerStateFn {
		line, err := l.peekLine()
		// Burn empty lines
		if isBlank(line) {
			err := l.immolateLine()
			if err != nil {
				l.emitError(err)
				return nil
			}
			return lexBlock
		}

		r, err := l.peek()
		if err != nil {
			l.emitError(err)
			return nil
		}

		// Burn and store the rune
		if isWhiteSpace(r) {
			return fn
		}

		if isParagraph(r) {
			return lexParagraph
		}
		return fn
	}

	return fn
}

func lexWorthlessSpace(l *blockLexer) lexerStateFn {
	l.immolateRune()
	return lexBlock
}

func lexBlockQuote(l *blockLexer) lexerStateFn {
	for {
		r, err := l.peek()
		if err != nil {
			break
		}

		if !isBlockQuote(r) {
			break
		}

		err = l.consumeLine()
		if err != nil {
			break
		}

	}

	if len(l.lines) > 0 {
		l.emit(BT_BLOCK_QUOTE)
	}

	return lexBlock
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
	for {
		r, err := l.peek()
		if err != nil {
			break
		}

		if !isParagraph(r) {
			break
		}

		err = l.consumeLine()
		if err != nil {
			break
		}

	}

	if len(l.lines) > 0 {
		l.emit(BT_PARAGRAPH)
	}

	return lexBlock
}

func isWhiteSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isParagraph(r rune) bool {
	// FIXME: this is totally fmted
	//        if a line it's indented
	//        w/ 4 spaces, it's a code block
	i := strings.IndexRune("*-+\n>", r)
	return i < 0
}

func isBlockQuote(r rune) bool {
	return r == '>'
}

func isBlank(line []byte) bool {
	return false
}
