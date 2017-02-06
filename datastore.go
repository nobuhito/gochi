package gochi

import (
	"fmt"
	"net/http"

	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
)

type Datastore struct {
	Goon    *goon.Goon
	Gochi   *Gochi
	Request *http.Request
}

func (g *Gochi) NewDatastore(r *http.Request) Datastore {
	return Datastore{
		Goon:    goon.NewGoon(r),
		Gochi:   g,
		Request: r,
	}
}

func (ds *Datastore) Put(data interface{}) error {
	_, err := ds.Goon.Put(data)
	return err
}

func (ds *Datastore) Del(data interface{}) error {
	key := ds.Goon.Key(data)
	return ds.Goon.Delete(key)
}

type SearchFilter struct {
	Key   string
	Value interface{}
}

type SearchQuery struct {
	Kind    string
	Filters []SearchFilter
	Order   string
	Limit   int
}

func (g *Gochi) NewSeachQuery(r *http.Request) SearchQuery {
	return SearchQuery{
		Kind:    "",
		Filters: nil,
		Order:   "",
		Limit:   100,
	}
}

func (sq *SearchQuery) GetQuery() (datastore.Query, error) {

	query := datastore.NewQuery("")
	if sq.Kind == "" {
		return *query, fmt.Errorf("%s", "kind is not set")
	}
	query = datastore.NewQuery(sq.Kind)

	if len(sq.Filters) > 0 {
		for _, v := range sq.Filters {
			query = query.Filter(v.Key, v.Value)
		}
	}

	if sq.Order != "" {
		query = query.Order(sq.Order)
	}

	if sq.Limit > 0 {
		query = query.Limit(sq.Limit)
	}

	return *query, nil
}
