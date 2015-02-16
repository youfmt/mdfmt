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

	ins := make([]os.FileInfo, len(files)/2)
	outs := make([]os.FileInfo, len(files)/2)

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

}
