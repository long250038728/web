package encode

import (
	"errors"
	"net/http"
	"net/url"
)

type fileWriter struct {
	writer
}

func (w *fileWriter) Run(response any, err error) {
	if err != nil {
		w.WriteErr(err)
		return
	}
	if file, ok := response.(File); !ok {
		w.WriteErr(errors.New("the File is not interface to File"))
	} else {
		w.WriteFile(file.FileName(), file.FileData())
	}
}

func (w *fileWriter) WriteErr(err error) {
	w.ginContext.JSON(http.StatusOK, NewResponse(nil, err))
}

//========================================================================================

func (w *fileWriter) WriteFile(fileName string, data []byte) {
	fileName = url.QueryEscape(fileName) // 防止中文乱码
	w.ginContext.Header("Content-Type", "application/octet-stream")
	w.ginContext.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	_, _ = w.ginContext.Writer.Write(data)
	w.addLog("File is write ok")
}
