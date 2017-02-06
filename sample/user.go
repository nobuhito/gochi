package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nobuhito/gochi"
)

func user() {
	user := g.Group("/api/user")
	user.Get("/:id", getUserByID)
	user.Put("/", putUser)
	user.Delete("/:id", deleteUser)
}

type User struct {
	ID   string `json:"id" datastore:"-" goon:"id"`
	Name string `json:"name"`
}

func getUserByID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	res := g.NewResponse(w)
	res.Body = []byte("Index")
	res.Write(w)
}

func putUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	users := Users{
		User{ID: "1", Name: "hogehoge"},
		User{ID: "2", Name: "hugahuga"},
	}

	ds := g.NewDatastore(r)
	for _, v := range users {
		err := ds.Put(&v)
		if err != nil {
			gochi.LogDebug(r, err)
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := User{ID: "1", Name: "hogehoge"}
	ds := g.NewDatastore(r)

	ds.Del(user)
}
