package gochi

import (
	"fmt"
	"reflect"

	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Datastore struct {
	Goon  *goon.Goon
	Gochi *Gochi
}

func (g *Gochi) NewDatastore(ctx context.Context) Datastore {
	return Datastore{
		Goon:  goon.FromContext(ctx),
		Gochi: g,
	}
}

func (ds *Datastore) GetNamespace(data interface{}) string {
	return fmt.Sprintf("%s", reflect.TypeOf(data).Name())
}

func (ds *Datastore) Put(data interface{}) error {
	_, err := ds.Goon.Put(data)
	return err
}

func (ds *Datastore) PutMulti(data interface{}) ([]*datastore.Key, error) {
	return ds.Goon.PutMulti(data)
}

func (ds *Datastore) Get(data interface{}) error {
	return ds.Goon.Get(data)
}

func (ds *Datastore) GetMulti(data interface{}) error {
	return ds.Goon.GetMulti(data)
}

func (ds *Datastore) GetAll(q *datastore.Query, data interface{}) ([]*datastore.Key, error) {
	return ds.Goon.GetAll(q, data)
}

func (ds *Datastore) GetBindAll(sq *SearchQuery, filters []SearchFilter, data interface{}) ([]*datastore.Key, error) {
	var keys []*datastore.Key
	for _, filter := range filters {
		var searchQuery SearchQuery
		searchQuery = sq.Clone()
		searchQuery.AppendFilter(filter)
		q, err := searchQuery.GetQuery()
		if err != nil {
			return nil, err
		}

		_keys, err := ds.Goon.GetAll(&q, data)
		if err != nil {
			return nil, err
		}

		for _, key := range _keys {
			keys = appendIfMissingKey(keys, key)
		}
	}
	return keys, nil
}

func appendIfMissingKey(keys []*datastore.Key, key *datastore.Key) []*datastore.Key {
	for _, _key := range keys {
		if _key.StringID() == key.StringID() {
			return keys
		}
	}
	return append(keys, key)
}

func (ds *Datastore) Del(data interface{}) error {
	key := ds.Goon.Key(data)
	return ds.Goon.Delete(key)
}

func (ds *Datastore) DelAll(namespace string) error {
	q := datastore.NewQuery(namespace).KeysOnly()
	for {
		keys, err := ds.Goon.GetAll(q, nil)
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return nil
		}
		err = ds.Goon.DeleteMulti(keys)
		if err != nil {
			return err
		}
	}
}

type SearchFilter struct {
	Key   string
	Value interface{}
}

type SearchQuery struct {
	Kind    string
	Filters []SearchFilter
	Orders  []string
	Limit   int
	Offset  int
}

func (g *Gochi) NewSeachQuery() *SearchQuery {
	return &SearchQuery{
		Filters: nil,
		Limit:   100,
		Offset:  0,
	}
}

func (g *Gochi) BuildSearchQuery(filters []SearchFilter) *SearchQuery {
	return &SearchQuery{
		Filters: filters,
		Limit:   1000,
		Offset:  0,
	}
}

func (sq *SearchQuery) Clone() SearchQuery {
	new := SearchQuery{}
	new.Kind = sq.Kind
	new.Filters = sq.Filters
	new.Orders = sq.Orders
	new.Limit = sq.Limit
	new.Offset = sq.Offset
	return new
}

func (sq *SearchQuery) SetKind(kind string) *SearchQuery {
	sq.Kind = kind
	return sq
}

func (sq *SearchQuery) SetLimit(limit int) *SearchQuery {
	sq.Limit = limit
	return sq
}

func (sq *SearchQuery) SetOffset(offset int) *SearchQuery {
	sq.Offset = offset
	return sq
}

func (sq *SearchQuery) SetOrder(order string) *SearchQuery {
	sq.Orders = append(sq.Orders, order)
	return sq
}

func (sq *SearchQuery) AppendFilter(filter SearchFilter) *SearchQuery {
	sq.Filters = append(sq.Filters, filter)
	return sq
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

	if len(sq.Orders) > 0 {
		for _, v := range sq.Orders {
			query = query.Order(v)
		}
	}

	if sq.Limit > 0 {
		query = query.Limit(sq.Limit)
	}

	if sq.Offset > 0 {
		query = query.Offset(sq.Offset)
	}

	return *query, nil
}
