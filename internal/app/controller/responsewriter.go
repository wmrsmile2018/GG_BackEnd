package controller

import "net/http"

type responseWrite struct {
	http.ResponseWriter // анонимное поле
	code                int
}

func (w *responseWrite) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
