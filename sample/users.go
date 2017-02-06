package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nobuhito/gochi"
)

func users() {
	users := g.Group("/api/users")
	users.Get("/", getUsers)
	users.Get("/:name", searchUsers)
}

type Users []User

func getUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	users := Users{
		User{ID: "1", Name: "hogehoge"},
		User{ID: "2", Name: "hugahuga"},
	}

	res := g.NewResponse(w)
	res.Body = users
	res.WriteJSON(w)
}

func searchUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sq := g.NewSeachQuery(r)
	sq.Kind = "User"
	filter := gochi.SearchFilter{
		Key:   "Name =",
		Value: p.ByName("name"),
	}
	sq.Filters = append(sq.Filters, filter)

	query, err := sq.GetQuery()
	if err != nil {
		gochi.LogDebug(r, err)
	}

	var users Users

	ds := g.NewDatastore(r)
	ds.Goon.GetAll(&query, &users)

}
