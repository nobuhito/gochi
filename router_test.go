package gochi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
)

func TestParamKeys(t *testing.T) {
	g := New()

	var tests = []struct {
		path string
		keys string
	}{
		{"/test?foo=abc&foo=def&key=foo", "foo, id, key"},
		{"/test/?foo=abc&key=foo", "foo, id, key"},
		{"/12345?key=id", "id, key"},
		{"/test/12345?key=foo", "foo, key"},
		{"/test/test.txt?key=foo", "foo, key"},
		{"/test/abc?foo=def&key=foo", "foo, key"},
		{"/test/abc/def?key=bar", "bar, foo, key"},
		{"/12345?key=invalidKey", "id, key"},
		{"/root/hoge?key=child", "child, key"},
		{"/test/test?FOO=test", "FOO, foo"},
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", paramKeysHandler)
	r.HandleFunc("/{id}", paramKeysHandler)
	r.HandleFunc("/test", paramKeysHandler)
	r.HandleFunc("/test/{foo}", paramKeysHandler)
	r.HandleFunc("/test/{foo}/{bar}", paramKeysHandler)

	r.PathPrefix("/root").Subrouter().HandleFunc("/{child}", paramKeysHandler).Methods("GET")

	server := httptest.NewServer(r)
	defer server.Close()

	for i, test := range tests {
		res, err := http.Get(server.URL + test.path)
		g.Ok(t, err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		g.Ok(t, err)

		exp := test.keys
		act := fmt.Sprintf("%s", body)
		g.Assert(t, exp == act, "\n %v)\t    exp: \"%+v\"\n\tbut got: \"%+v\"", i+1, exp, act)
	}

}

func TestVars(t *testing.T) {

	g := New()

	var tests = []struct {
		path string
		key  string
		val  string
	}{
		{"/test?foo=abc&foo=def&key=foo", "foo", "abc"},
		{"/test/?foo=abc&key=foo", "foo", "abc"},
		{"/12345?key=id", "id", "12345"},
		{"/test/12345?key=foo", "foo", "12345"},
		{"/test/test.txt?key=foo", "foo", "test.txt"},
		{"/test/abc?foo=def&key=foo", "foo", "abc"},
		{"/test/abc/def?key=bar", "bar", "def"},
		{"/12345?key=invalidKey", "invalidKey", ""},
		{"/root/hoge?key=child", "child", "hoge"},
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", nilHandler)
	r.HandleFunc("/{id}", varsHandler)
	r.HandleFunc("/test", varsHandler)
	r.HandleFunc("/test/{foo}", varsHandler)
	r.HandleFunc("/test/{foo}/{bar}", varsHandler)

	r.PathPrefix("/root").Subrouter().HandleFunc("/{child}", varsHandler).Methods("GET")

	server := httptest.NewServer(r)
	defer server.Close()

	for i, test := range tests {

		res, err := http.Get(server.URL + test.path)
		g.Ok(t, err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		g.Ok(t, err)

		exp := test.val
		act := fmt.Sprintf("%s", body)
		g.Assert(t, exp == act, "\n %v)\t    exp: \"%+v\"\n\tbut got: \"%+v\"", i+1, exp, act)
	}
}

func paramKeysHandler(w http.ResponseWriter, r *http.Request) {
	g := New()
	keys := g.ParamKeys(r)
	sort.Strings(keys)
	w.Write([]byte(strings.Join(keys, ", ")))
}

func nilHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(nil)
}

func varsHandler(w http.ResponseWriter, r *http.Request) {
	g := New()
	key := g.Vars(r, "key")
	val := g.Vars(r, key)
	w.Write([]byte(val))
}

func TestStatic(t *testing.T) {
	g := New()

	server := httptest.NewServer(g.Router)
	defer server.Close()

	var tests = []struct {
		path string
		out  string
	}{
		{"/index.html", "Hello Gochi"},
		{"/static/js/script.js", "console.log(\"Hello Gochi\");"},
		{"/notfound.html", "not found: notfound.html"},
		{"/foo/bar", "404 page not found"},
	}

	for _, test := range tests {
		res, err := http.Get(server.URL + test.path)
		g.Ok(t, err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		g.Ok(t, err)

		ret := strings.TrimSpace(fmt.Sprintf("%s", body))
		g.Equals(t, test.out, ret)
	}
}

func TestGroupPath(t *testing.T) {
	g := New()

	inst, _, spindown := g.SpinUp(t)
	defer spindown()

	req, err := inst.NewRequest("GET", "/api/test", nil)
	g.Ok(t, err)

	w := httptest.NewRecorder()

	grp := g.Group("/api")
	grp.GET("/test", handler)
	g.Router.ServeHTTP(w, req)
	g.Equals(t, "{\"id\":\"abcdefg\"}", w.Body.String())
}

func handler(ctx context.Context, r *http.Request) Response {
	success := struct {
		ID string `json:"id"`
	}{ID: "abcdefg"}

	res := ResponseJSON(http.StatusOK, success)
	return res
}

func TestJSONWriter(t *testing.T) {
	g := New()

	inst, _, spindown := g.SpinUp(t)
	defer spindown()

	req, err := inst.NewRequest("GET", "/test", nil)
	g.Ok(t, err)

	w := httptest.NewRecorder()

	g.GET("/test", handler)
	g.Router.ServeHTTP(w, req)
	g.Equals(t, "{\"id\":\"abcdefg\"}", w.Body.String())
}

func TestMakeRes(t *testing.T) {
	g := New()

	inst, _, spindown := g.SpinUp(t)
	defer spindown()

	req, err := inst.NewRequest("GET", "/test", nil)
	g.Ok(t, err)

	w := httptest.NewRecorder()

	g.GET("/test", handlerJSON)
	g.Router.ServeHTTP(w, req)
	_, out := testJSON()
	g.Equals(t, out, w.Body.String())
}

func handlerJSON(ctx context.Context, r *http.Request) Response {
	response := ResponseJSON(http.StatusOK, retJSON())
	return response
}
