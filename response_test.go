package gochi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespond(t *testing.T) {
	g := New()
	status := http.StatusOK
	bodyString := "test"
	bodyJSON, out := testJSON()

	responseString := respond(status, bodyString)
	g.Equals(t, status, responseString.status)
	g.Equals(t, bodyString, string(responseString.body))

	responseJSON := respond(status, bodyJSON)
	g.Equals(t, status, responseJSON.status)
	g.Equals(t, out, string(responseJSON.body))
}

func TestResponseError(t *testing.T) {
	g := New()
	w := httptest.NewRecorder()
	response := ResponseErrorJSON(http.StatusInternalServerError, "error")
	response.Write(w)
	g.Equals(t, "error", w.Body.String())
}

func TestResponseCreated(t *testing.T) {
	g := New()
	w := httptest.NewRecorder()
	body, out := testJSON()
	response := ResponseCreated(http.StatusCreated, body, "")
	response.Write(w)
	g.Equals(t, out, w.Body.String())
}

func TestResponseJSON(t *testing.T) {
	g := New()
	w := httptest.NewRecorder()
	body, out := testJSON()
	response := ResponseJSON(http.StatusOK, body)
	response.Write(w)
	g.Equals(t, out, w.Body.String())
}

func TestResponseEmpty(t *testing.T) {
	g := New()
	w := httptest.NewRecorder()
	response := ResponseEmpty(http.StatusNoContent)
	response.Write(w)
	g.Equals(t, "null", w.Body.String())
}
