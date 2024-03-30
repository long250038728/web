package gen

import (
	"bufio"
	"bytes"
	"go/format"
	"io"
	"os"
	"text/template"
)

type GenImpl struct {
	Name     string
	TmplPath string
	Func     template.FuncMap
	Data     any
	IsFormat bool
}

func (g *GenImpl) Gen() ([]byte, error) {
	var file *os.File
	var contents []byte
	var tmpl *template.Template
	var err error

	//读取文件字段
	if file, err = os.Open(g.TmplPath); err != nil {
		return nil, err
	}
	if contents, err = io.ReadAll(file); err != nil {
		return nil, err
	}

	//生成tmpl
	buffer := new(bytes.Buffer)
	writer := bufio.NewWriter(buffer)
	if tmpl, err = template.New(g.Name).Funcs(g.Func).Parse(string(contents)); err != nil {
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
