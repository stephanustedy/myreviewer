package util

import(
	"bytes"
	"net/http"
)

func RenderHTML(w http.ResponseWriter, buffer bytes.Buffer, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buffer.Bytes())
}