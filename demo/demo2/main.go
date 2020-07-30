package main

import (
	"bytes"
	"io/ioutil"
	"path"
	"runtime"

	ext "github.com/OhYee/goldmark-fenced_codeblock_extension"
	python "github.com/OhYee/goldmark-python"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var raw = "```python\nimport matplotlib.pyplot as plt\nimport matplotlib\nimport sys\n\nmatplotlib.use(\"svg\")\n\ny = [1, 2, 3, 4, 5]\nx = [5, 4, 3, 2, 1]\n\nplt.plot(x, y)\nplt.savefig(sys.stdout)\n```\n\n```python-output\nimport matplotlib.pyplot as plt\nimport matplotlib\nimport sys\n\nmatplotlib.use(\"svg\")\n\ny = [1, 2, 3, 4, 5]\nx = [5, 4, 3, 2, 1]\n\nplt.plot(x, y)\nplt.savefig(sys.stdout)\n```"

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			ext.NewExt(
				python.RenderMap(20, "python3", "python-output"),
				ext.RenderMap{
					Languages:      []string{"*"},
					RenderFunction: ext.GetFencedCodeBlockRendererFunc(highlighting.NewHTMLRenderer()),
				},
			),
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
