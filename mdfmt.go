package mdfmt

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func Fmt(input io.Reader) (output io.Reader) {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	out := bytes.NewBuffer(make([]byte, 0, 1024))

	for scanner.Scan() {
		line := scanner.Text()
		out.WriteString(line + "\n")

		switch {
		case line == "":
		case isHeader(line):
			mustBeFollowedByBlankLine(scanner, out)
		}

	}

	return out
}

func isHeader(line string) bool {
	return line[0] == '#'
}

func mustBeFollowedByBlankLine(scanner *bufio.Scanner, out io.Writer) {
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" {
			out.Write([]byte("\n"))
		} else {
			out.Write([]byte("\n"))
			out.Write([]byte(scanner.Text() + "\n"))
		}
	}

	return
}
