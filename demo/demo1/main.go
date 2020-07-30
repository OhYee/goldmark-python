package main

import (
	"bytes"
	"io/ioutil"
	"path"
	"runtime"

	python "github.com/OhYee/goldmark-python"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var raw = "```python\nprint(\"Hello World\")"

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			python.Default,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(),
	)
	buf := bytes.Buffer{}
	if err := md.Convert([]byte(raw), &buf); err != nil {
		panic(err.Error())
	}

	_, file, _, _ := runtime.Caller(0)
	if err := ioutil.WriteFile(path.Join(path.Dir(file), "output.html"), buf.Bytes(), 777); err != nil {
		panic(err.Error())
	}
}
