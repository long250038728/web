package gateway

import (
	"errors"
	"net/url"
)

type FileInterface interface {
	FileName() string
	FileData() []byte
}

type fileWriter struct {
	writer
}

func (w *fileWriter) Run(response any, err error) {
	if err != nil {
		w.WriteErr(err)
		return
	}
	if file, ok := response.(FileInterface); !ok {
		w.WriteErr(errors.New("the File is not interface to FileInterface"))
	} else {
		w.WriteFile(file.FileName(), file.FileData())
	}

}
func (w *writer) WriteFile(fileName string, data []byte) {
	fileName = url.QueryEscape(fileName) // 防止中文乱码
	w.ginContext.Header("Content-Type", "application/octet-stream")
	w.ginContext.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	_, _ = w.ginContext.Writer.Write(data)
	w.addLog("File is write ok")
}
