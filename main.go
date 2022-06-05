package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const header = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title> Yet Another Markdown Preview </title>
	</head>
	<body>`

const footer = `
	</body>
</html>`

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	skip := flag.Bool("skip", false, "Skip auto-preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skip); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, out io.Writer, skipPreview bool) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(file)

	tempfile, err := ioutil.TempFile("", "mdp_"+filepath.Base(filename)+"_*.html")
	if err != nil {
		return err
	}
	if err := tempfile.Close(); err != nil {
		return err
	}
	genFile := tempfile.Name()
	fmt.Fprintf(out, "Saving preview html file in %s\n", genFile)
	if err := saveHtml(genFile, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}
	defer os.Remove(genFile)
	if err := preview(genFile); err != nil {
		return nil
	}
	time.Sleep(2 * time.Second)
	return nil
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

func preview(filename string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("os is not supported")
	}

	cParams = append(cParams, filename)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	return exec.Command(cPath, cParams...).Run()
}
