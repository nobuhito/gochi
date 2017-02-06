package main

import (
	"net/http"

	"github.com/nobuhito/gochi"
)

var g *gochi.Gochi

func init() {

	g = gochi.New()
	g.InDev = true
	g.Static("public")

	user()
	users()
	contents()

	http.Handle("/", g.Router)
}
