package gochi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/luci/luci-go/common/logging"
)

type Response interface {
	Write(w http.ResponseWriter)
	Status() int
}

type baseResponse struct {
	status int
	body   []byte
	header http.Header
}

func (r *baseResponse) Write(w http.ResponseWriter) {
	header := w.Header()
	for k, v := range r.header {
		header[k] = v
	}
	w.WriteHeader(r.status)
	w.Write(r.body)
}

func (r *baseResponse) Status() int {
	return r.status
}

func (r *baseResponse) Header(key, value string) *baseResponse {
	r.header.Set(key, value)
	return r
}

func ResponseEmpty(status int) *baseResponse {
	return respond(status, nil)
}

func ResponseJSON(status int, body interface{}) *baseResponse {
	return respond(status, body).Header("Content-Type", "application/json; charset=utf8")
}

func ResponseCreated(status int, body interface{}, location string) *baseResponse {
	return ResponseJSON(status, body).Header("Location", location)
}

func ResponseErrorJSON(status int, message string) *baseResponse {
	return respond(status, message).Header("Content-Type", "application/json; charset=utf8")
}

func respond(status int, body interface{}) *baseResponse {
	var b []byte
	var err error

	switch t := body.(type) {
	case string:
		b = []byte(t)
	default:
		b, err = json.Marshal(body)
		if err != nil {
			return ResponseErrorJSON(
				http.StatusInternalServerError,
				fmt.Sprintf("faild marshalling json: %s", err.Error()),
			)
		}
	}

	return &baseResponse{
		status: status,
		body:   b,
		header: make(http.Header),
	}
}

func (g *Gochi) IfErrReturnJsonResponse(ctx context.Context, err error) Response {
	if err != nil {
		logging.Errorf(ctx, "%+v", err)
		return ResponseJSON(http.StatusInternalServerError, err.Error())
	}
	return nil
}
