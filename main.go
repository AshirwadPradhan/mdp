package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const header = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type content="text/html; charset=utf-8">
		<title> Yet Another Markdown Preview </title>
	</head>
	<body>`

const footer = `
	</body>
</html>`

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(file)

	previewOut := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Printf("Saving preview html file in %s\n", previewOut)

	return saveHtml(previewOut, htmlData)
}

func parseContent(markdown []byte) []byte {
	convertedHTML := blackfriday.Run(markdown)
	sanitizedHTML := bluemonday.UGCPolicy().SanitizeBytes(convertedHTML)

	var parsedHTML bytes.Buffer

	parsedHTML.WriteString(header)
	parsedHTML.WriteString(string(sanitizedHTML))
	parsedHTML.WriteString(footer)

	return parsedHTML.Bytes()
}

func saveHtml(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}
