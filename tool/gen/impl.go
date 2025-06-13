package gen

import (
	"bufio"
	"bytes"
	"go/format"
	"io"
	"os"
	"text/template"
)

type Impl struct {
	Name        string
	TmplPath    string
	Tmpl        string
	Func        template.FuncMap
	Data        any
	IsFormat    bool
	delimsLeft  string
	delimsRight string
}

func (g *Impl) SetDelims(delimsLeft, delimsRight string) *Impl {
	g.delimsLeft = delimsLeft
	g.delimsRight = delimsRight
	return g
}

func (g *Impl) Gen() ([]byte, error) {
	var file *os.File
	var contents []byte
	var tmpl *template.Template
	var err error

	if len(g.Tmpl) > 0 {
		contents = []byte(g.Tmpl)
	}

	if len(g.TmplPath) > 0 {
		//读取文件字段
		if file, err = os.Open(g.TmplPath); err != nil {
			return nil, err
		}
		if contents, err = io.ReadAll(file); err != nil {
			return nil, err
		}
	}

	//生成tmpl
	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)

	tmpl = template.New(g.Name)

	if len(g.delimsLeft) > 0 && len(g.delimsRight) > 0 {
		tmpl = tmpl.Delims(g.delimsLeft, g.delimsRight)
	}

	if tmpl, err = tmpl.Funcs(g.Func).Parse(string(contents)); err != nil {
		return nil, err
	}
	if err := tmpl.Execute(writer, g.Data); err != nil {
		return nil, err
	}
	if err = writer.Flush(); err != nil {
		return nil, err
	}

	//是否format go语言的内容
	if g.IsFormat {
		return format.Source(buffer.Bytes())
	}
	return buffer.Bytes(), nil
}
