package gochi

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// simple route

func (g *Gochi) Static(path string) {
	g.Router.NotFound = http.FileServer(http.Dir(path))
}

func (g *Gochi) Get(path string, handle httprouter.Handle) {
	g.Router.GET(path, handle)
}

func (g *Gochi) Put(path string, handle httprouter.Handle) {
	g.Router.PUT(path, handle)
}

func (g *Gochi) Delete(path string, handle httprouter.Handle) {
	g.Router.DELETE(path, handle)
}

func (g *Gochi) Post(path string, handle httprouter.Handle) {
	g.Router.POST(path, handle)
}

// group route

type Group struct {
	Group string
	Gochi *Gochi
}

func (g *Gochi) Group(root string) Group {
	return Group{
		Group: root,
		Gochi: g,
	}
}

func (grp *Group) Get(path string, handle httprouter.Handle) {
	grp.Gochi.Get(grp.Group+path, handle)
}

func (grp *Group) Put(path string, handle httprouter.Handle) {
	grp.Gochi.Put(grp.Group+path, handle)
}

func (grp *Group) Delete(path string, handle httprouter.Handle) {
	grp.Gochi.Delete(grp.Group+path, handle)
}

func (grp *Group) Post(path string, handle httprouter.Handle) {
	grp.Gochi.Post(grp.Group+path, handle)
}
