package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test_basic.md"
	goldenFile = "./testdata/test_basic.md.html"
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

	var mockStdOut bytes.Buffer
	if err := run(inputFile, &mockStdOut); err != nil {
		t.Fatal(err)
	}
	parseStdout := strings.Split(strings.TrimSpace(mockStdOut.String()), " ")
	outFile := parseStdout[len(parseStdout)-1]
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
