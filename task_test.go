package gochi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"google.golang.org/appengine/taskqueue"
)

// var g *Gochi

type Data struct {
	ID    string
	value string
}

var g_ *Gochi // for TestDelayGet

func TestDelayGet(t *testing.T) {
	g_ = New()

	var tests = []struct {
		path string
		f    interface{}
		body string
	}{
		{"/success", successFunc, "null"},
		{"/error", errorFunc, "null"},
	}

	inst, ctx, spindown := SpinUp(t)
	defer spindown()

	for _, test := range tests {
		g_.DelayGET(test.path, testHandler, test.f)
	}

	http.Handle("/", g_.Router)

	for _, test := range tests {

		r, err := inst.NewRequest("GET", test.path, nil)
		g_.Ok(t, err)

		response := testHandler(ctx, r)
		w := httptest.NewRecorder()
		response.Write(w)
		g_.Equals(t, test.body, w.Body.String())
		// g.Equals(t, 200, response.Status())
	}

	q, err := taskqueue.QueueStats(ctx, []string{"default"})
	g_.Ok(t, err)
	g_.Equals(t, 2, q[0].Tasks)
}

func testHandler(ctx context.Context, r *http.Request) Response {
	// ctx := appengine.NewContext(r)
	task, err := g_.SearchTask(r)
	if err != nil {
		return ResponseErrorJSON(http.StatusInternalServerError, err.Error())
	}
	err = task.Call(ctx)
	if err != nil {
		return ResponseErrorJSON(http.StatusInternalServerError, fmt.Sprintf("%+v", err))
	}

	return ResponseEmpty(http.StatusOK)
}

var successFunc = func(ctx context.Context) error {
	return nil
}

var errorFunc = func(ctx context.Context) error {
	// TODO: 実行結果までテスト
	return errors.New("error")
}
