package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

const (
	inputFile  = "./testdata/test_basic.md"
	goldenFile = "./testdata/test_basic.md.html"
	outFile    = "test_basic.md.html"
)

func TestParseContent(t *testing.T) {
	file, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	parsedHTML := parseContent(file)

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, parsedHTML) {
		t.Logf("golden: \n%s\n", expected)
		t.Logf("genfile: \n%s\n", parsedHTML)
		t.Error("Result does not matches golden file")
	}

}

func TestRun(t *testing.T) {
	if err := run(inputFile); err != nil {
		t.Fatal(err)
	}

	gen, err := ioutil.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}

	exp, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(gen, exp) {
		t.Logf("golden: \n%s\n", exp)
		t.Logf("genfile: \n%s\n", gen)
		t.Error("Result does not matches golden file")
	}
	os.Remove(outFile)
}
