package mdfmt

import (
	"bufio"
	"bytes"
	"io"
)

func Fmt(input io.Reader) (output io.Reader) {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	out := bytes.NewBuffer(make([]byte, 0, 1024))

	for scanner.Scan() {
		line := scanner.Text()

		if line[0] == '#' {
			out.WriteString(line + "\n")
			mustBeFollowedByBlankLine(scanner, out)
		}
	}

	return out
}

func mustBeFollowedByBlankLine(scanner *bufio.Scanner, out io.Writer) {
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		if scanner.Text() == "" {
			out.Write([]byte("\n"))
		} else {
			out.Write([]byte("\n"))
			out.Write([]byte(scanner.Text() + "\n"))
		}
	}

	return
}
