package gochi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

type Group struct {
	Parent string
	Gochi  *Gochi
}

func initRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(`/{file:[^/.]*\.html}`, staticHandler)
	router.HandleFunc("/", indexHtmlHandler)

	static := router.PathPrefix("/static/")
	static.Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	return router
}

// type ParamBy func(key string) (value string)

func (g *Gochi) ParamKeys(r *http.Request) (keys []string) {
	queries := r.URL.Query()
	for key, _ := range queries {
		keys = append(keys, key)
	}

	vars := mux.Vars(r)
	for key, _ := range vars {
		keys = append(keys, key)
	}

	results := make([]string, 0, len(keys))
	encountered := map[string]bool{}
	for _, key := range keys {
		if !encountered[key] {
			encountered[key] = true
			results = append(results, key)
		}
	}

	return results
}

func (g *Gochi) Vars(r *http.Request, key string) string {
	vars := mux.Vars(r)
	ret, ok := vars[key]
	if ok {
		return ret
	}

	query := r.URL.Query()

	ret = query.Get(key)
	return ret
}

func responseToHandler(g *Gochi, h func(r *http.Request) Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := h(r)
		response.Write(w)
	}
}

func (g *Gochi) GET(path string, h func(r *http.Request) Response) {
	g.Router.HandleFunc(path, responseToHandler(g, h)).Methods("GET")
}

func (g *Gochi) PUT(path string, h func(r *http.Request) Response) {
	g.Router.HandleFunc(path, responseToHandler(g, h)).Methods("PUT")
}

func (g *Gochi) POST(path string, h func(r *http.Request) Response) {
	g.Router.HandleFunc(path, responseToHandler(g, h)).Methods("POST")
}

func (g *Gochi) DELETE(path string, h func(r *http.Request) Response) {
	g.Router.HandleFunc(path, responseToHandler(g, h)).Methods("DELETE")
}

func (g *Gochi) Group(path string) Group {
	group := Group{Parent: path, Gochi: g}
	return group
}

func (g *Group) GET(path string, h func(r *http.Request) Response) {
	g.Gochi.Router.PathPrefix(g.Parent).Subrouter().HandleFunc(path, responseToHandler(g.Gochi, h)).Methods("GET")
}

func (g *Group) PUT(path string, h func(r *http.Request) Response) {
	g.Gochi.Router.PathPrefix(g.Parent).Subrouter().HandleFunc(path, responseToHandler(g.Gochi, h)).Methods("PUT")
}

func (g *Group) POST(path string, h func(r *http.Request) Response) {
	g.Gochi.Router.PathPrefix(g.Parent).Subrouter().HandleFunc(path, responseToHandler(g.Gochi, h)).Methods("POST")
}

func (g *Group) DELETE(path string, h func(r *http.Request) Response) {
	g.Gochi.Router.PathPrefix(g.Parent).Subrouter().HandleFunc(path, responseToHandler(g.Gochi, h)).Methods("DELETE")
}

func indexHtmlHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["file"]
	path := path.Join("public", file)
	_, err := os.Stat(path)
	if err == nil {
		// tmpl := template.Must(template.ParseFiles("public/" + file))
		// err := tmpl.Execute(w, nil)
		bs, err := ioutil.ReadFile("public/" + file)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		} else {
			fmt.Fprintf(w, "%s", string(bs))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found: %s", file)
	}
}
