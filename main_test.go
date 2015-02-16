package mdfmt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testCasesPath string = "./test_cases/"
)

func TestFmt(t *testing.T) {
	files, err := ioutil.ReadDir(testCasesPath)
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

	for i, in := range ins {
		inFile, err := os.Open(filepath.Join(testCasesPath, in.Name()))
		if err != nil {
			t.Fatal(err)
		}
		defer inFile.Close()

		outFile, err := os.Open(filepath.Join(testCasesPath, outs[i].Name()))
		if err != nil {
			t.Fatal(err)
		}
		defer outFile.Close()

		actualOutput, err := ioutil.ReadAll(Fmt(inFile))
		if err != nil {
			t.Fatal(err)
		}

		expectedOutput, err := ioutil.ReadAll(outFile)
		if err != nil {
			t.Fatal(err)
		}

		if string(actualOutput) != string(expectedOutput) {
			t.Logf("\nEXPECTED:\n%s\nGOT:\n%s\n", expectedOutput, actualOutput)
			t.Fail()
		}

	}

}
