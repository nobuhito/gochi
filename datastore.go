package gochi

import (
	"reflect"

	"go.chromium.org/gae/filter/dscache"
	"go.chromium.org/gae/service/datastore"
	"golang.org/x/net/context"
)

type Datastore struct {
	Context     context.Context
	UseMemcache bool
}

var ErrNoSuchEntity = datastore.ErrNoSuchEntity

func (g *Gochi) NewDatastore(ctx context.Context) (Datastore, error) {

	testable := datastore.GetTestable(ctx)
	if testable != nil {
		testable.Consistent(true)
		testable.AutoIndex(true)
	}

	ds := Datastore{
		Context:     ctx,
		UseMemcache: true,
	}

	return ds, nil
}

func (ds *Datastore) useMemcache(on bool) {
	ds.UseMemcache = on
}

func (ds *Datastore) Put(data interface{}) error {
	ctx := ds.Context
	if ds.UseMemcache {
		ctx = dscache.FilterRDS(ctx)
	}
	return datastore.Put(ctx, data)
}

func (ds *Datastore) Get(data interface{}) error {
	ctx := ds.Context
	if ds.UseMemcache {
		ctx = dscache.FilterRDS(ctx)
	}
	return datastore.Get(ctx, data)
}

func (ds *Datastore) GetAll(q *datastore.Query, data interface{}) error {
	ctx := ds.Context
	if ds.UseMemcache {
		ctx = dscache.FilterRDS(ctx)
	}

	err := datastore.GetAll(ctx, q, data)
	return err
}

func (ds *Datastore) Delete(data interface{}) error {
	ctx := ds.Context
	if ds.UseMemcache {
		ctx = dscache.FilterRDS(ctx)
	}
	return datastore.Delete(ctx, data)
}

func (ds *Datastore) Run(query *datastore.Query, cb interface{}) error { // TODO: require test
	ctx := ds.Context
	if ds.UseMemcache {
		ctx = dscache.FilterRDS(ctx)
	}
	return datastore.Run(ctx, query, cb)
}

func (ds *Datastore) NewQuery(kind interface{}) *datastore.Query { // TODO: require test
	_kind, ok := kind.(string)
	if ok {
		return datastore.NewQuery(_kind)
	}
	return datastore.NewQuery(reflect.TypeOf(kind).Name())
}

type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

type QueryVars struct {
	Filter []Filter
	Limit  int32
	Offset int32
	Order  []string
}

func (ds *Datastore) BuildQuery(kind interface{}, queryVars QueryVars) (*datastore.Query, error) {
	q := ds.NewQuery(kind)

	if queryVars.Limit != 0 {
		q = q.Limit(queryVars.Limit).Offset(queryVars.Offset)
	}

	if len(queryVars.Order) > 0 {
		q = q.Order(queryVars.Order...)
	}

	for _, filter := range queryVars.Filter {
		f := filter.Field
		v := filter.Value
		switch filter.Operator {
		case "=":
			q = q.Eq(f, v)
		case ">":
			q = q.Lt(f, v)
		case ">=", "=>":
			q = q.Lte(f, v)
		case "<":
			q = q.Gt(f, v)
		case "<=", "=<":
			q = q.Gte(f, v)
		}
	}

	return q, nil
}
