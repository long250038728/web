package gateway

type jsonWriter struct {
	writer
}

func (w *jsonWriter) Run(response any, err error) {
	w.WriteJSON(response, err)
}
