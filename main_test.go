package mdfmt

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFmt(t *testing.T) {
	files, err := ioutil.ReadDir("./test_cases/")
	if err != nil {
		t.Fatal(err)
	}

	ins := make([]os.FileInfo, 0, len(files)/2)
	outs := make([]os.FileInfo, 0, len(files)/2)

	for _, file := range files {
		if strings.Contains(file.Name(), ".in.") {
			ins = append(ins, file)
		} else if strings.Contains(file.Name(), ".out.") {
			outs = append(outs, file)
		}
	}

	if len(ins) != len(outs) {
		t.Fatal("There is a different number of ins and out")
	}

	if len(ins) == 0 {
		t.Fatal("There are no input files")
	}

	if len(outs) == 0 {
		t.Fatal("There are no output files")
	}
}
