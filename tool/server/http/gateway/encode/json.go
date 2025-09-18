package encode

import (
	"encoding/json"
	"net/http"
)

type jsonWriter struct {
	writer
}

func (w *jsonWriter) Run(response any, err error) {
	w.WriteJSON(response, err)
}

func (w *jsonWriter) WriteErr(err error) {
	w.WriteJSON(nil, err)
}

//===================================================================

func (w *jsonWriter) WriteJSON(data interface{}, err error) {
	var res *Response
	var marshalByte []byte

	switch data.(type) {
	case Response:
		val, _ := data.(Response)
		res = &val
	case *Response:
		res = data.(*Response)
	default:
		res = NewResponse(data, err)
	}

	marshalByte, err = json.Marshal(res)
	if err != nil {
		res.Code = "999999"
		res.Message = err.Error()
		res.Details = nil
		marshalByte, _ = json.Marshal(res)
	}
	w.addLog(string(marshalByte))
	w.ginContext.JSON(http.StatusOK, res)
}

//func (w *jsonWriter) writeBytes(b []byte) {
//	//记录请求响应
//	w.addLog(string(b))
//
//	w.ginContext.Header("Content-Type", "application/JSON")
//	_, _ = w.ginContext.Writer.Write(b)
//}
