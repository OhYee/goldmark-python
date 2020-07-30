// Package goldmark-python is a extension for the goldmark(http://github.com/yuin/goldmark).
//
// This extension can run python in markdown

package python

import (
	"bytes"
	"crypto/sha1"
	"io"
	"os/exec"

	ext "github.com/OhYee/goldmark-fenced_codeblock_extension"
	fp "github.com/OhYee/goutils/functional"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// Default Python extension when there is no other fencedCodeBlock goldmark render extensions
var Default = NewPythonExtension(20, "python3", "python")

// RenderMap return the goldmark-fenced_codeblock_extension.RenderMap
func RenderMap(length int, pythonPath string, languages ...string) ext.RenderMap {
	return ext.RenderMap{
		Languages:      languages,
		RenderFunction: NewPython(length, pythonPath, languages...).Renderer,
	}
}

// NewPythonExtension return the goldmark.Extender
func NewPythonExtension(length int, pythonPath string, languages ...string) goldmark.Extender {
	return ext.NewExt(RenderMap(length, pythonPath, languages...))
}

// Python render struct
type Python struct {
	Languages  []string
	buf        map[string][]byte
	MaxLength  int
	PythonPath string
}

// NewPython initial a Python struct
func NewPython(length int, pythonPath string, languages ...string) *Python {
	return &Python{Languages: languages, buf: make(map[string][]byte), MaxLength: length, PythonPath: pythonPath}
}

// Renderer render function
func (d *Python) Renderer(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	language := string(n.Language(source))

	if fp.AnyString(func(l string) bool {
		return l == language
	}, d.Languages) {
		if !entering {
			raw := d.getLines(source, node)
			h := sha1.New()
			h.Write(raw)
			hash := string(h.Sum([]byte{}))
			if result, exist := d.buf[hash]; exist {
				w.Write([]byte(result))
			} else {
				res, err := runPython(raw, d.PythonPath)
				if err != nil {
					buf := bytes.NewBufferString("<p class=\"goldmark-python-error\">")
					buf.Write(res)
					buf.WriteString("</p>")
					res = buf.Bytes()
				}
				if len(d.buf) >= d.MaxLength {
					d.buf = make(map[string][]byte)
				}
				d.buf[hash] = res
				w.Write(res)
			}
		}
	}
	return ast.WalkContinue, nil
}

func (d *Python) getLines(source []byte, n ast.Node) []byte {
	buf := bytes.NewBuffer([]byte{})
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		buf.Write(line.Value(source))
	}
	return buf.Bytes()
}

func runPython(input []byte, pythonPath string, args ...string) (output []byte, err error) {
	var stdin io.WriteCloser

	pythonPath, err = exec.LookPath(pythonPath)
	if err != nil {
		return
	}

	cmd := exec.Command(pythonPath, args...)
	stdin, err = cmd.StdinPipe()
	if err != nil {
		return
	}

	_, err = stdin.Write(input)
	if err != nil {
		return
	}

	stdin.Close()
	if err != nil {
		return
	}

	output, err = cmd.CombinedOutput()
	if err != nil {
		return
	}

	return
}
