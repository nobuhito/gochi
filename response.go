package gochi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status int
	Body   interface{}
	Header http.Header
	InDev  bool
}

func (g *Gochi) NewResponse(w http.ResponseWriter) Response {
	return Response{
		Status: http.StatusOK,
		Body:   nil,
		Header: w.Header(),
		InDev:  g.InDev,
	}
}

func (res *Response) Write(w http.ResponseWriter) {
	header := w.Header()
	for k, v := range res.Header {
		header[k] = v
	}
	w.WriteHeader(res.Status)
	fmt.Fprintf(w, "%s", res.Body)
}

func (res *Response) WriteJSON(w http.ResponseWriter) {
	if res.InDev {
		res.Header.Set("Content-Type", "text/json: charset=UTF-8")
	} else {
		res.Header.Set("Content-Type", "application/json: charset=UTF-8")
	}

	header := w.Header()
	for k, v := range res.Header {
		header[k] = v
	}
	w.WriteHeader(res.Status)

	json.NewEncoder(w).Encode(res.Body)
}
