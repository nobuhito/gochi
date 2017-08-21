// +build !appengine

package gochi

import (
	"context"
	"strconv"
	"testing"

	"go.chromium.org/gae/impl/memory"
	"go.chromium.org/gae/service/datastore"
)

type Hogehoge struct {
	ID   string `goon:"id" gae:"$id"`
	Abcd string
	Efg  int
}
type Hogehoges []Hogehoge

func hoge() Hogehoge {
	return Hogehoge{
		ID:   "abcdefg",
		Abcd: "hogehoge is hoge",
		Efg:  1234,
	}
}

func TestDatastoreCRUD(t *testing.T) {
	t.Parallel()

	g := New()

	ctx := memory.Use(context.Background())
	ds, err := g.NewDatastore(ctx)
	g.Ok(t, err)

	hoge := hoge()
	g.Ok(t, ds.Put(&hoge))

	hogehoge := Hogehoge{
		ID: "abcdefg",
	}

	g.Ok(t, ds.Get(&hogehoge))
	g.Equals(t, hoge, hogehoge)

	g.Ok(t, ds.Delete(&hoge))

	e := "datastore: no such entity"
	g.Equals(t, e, ds.Get(&hoge).Error())
}

func TestUseMemcache(t *testing.T) {
	// TODO: 実装
}

func TestNewQuery(t *testing.T) {
	t.Parallel()

	g := New()

	ctx := memory.Use(context.Background())
	ds, err := g.NewDatastore(ctx)
	g.Ok(t, err)

	exp := ds.NewQuery("Hogehoge")
	act := ds.NewQuery(hoge())
	g.Equals(t, exp, act)
}

func TestPutMulti(t *testing.T) {
	t.Parallel()

	g := New()

	ctx := memory.Use(context.Background())
	ds, err := g.NewDatastore(ctx)
	g.Ok(t, err)

	hogehoges := []Hogehoge{}
	count := 1000
	for i := 0; i < count; i++ {
		hoge := hoge()
		hoge.ID = "id" + strconv.Itoa(i)
		hogehoges = append(hogehoges, hoge)
	}

	g.Ok(t, ds.Put(hogehoges))
}

func TestBuildQuery(t *testing.T) {
	t.Parallel()

	g := New()

	ctx := memory.Use(context.Background())
	ds, err := g.NewDatastore(ctx)
	g.Ok(t, err)

	query := ds.NewQuery("test")

	var tests = []struct {
		qv  QueryVars
		act *datastore.Query
	}{
		{qv: QueryVars{Limit: 10, Offset: 5, Order: []string{"-Create"}}, act: query.Limit(10).Offset(5).Order("-Create")},
		{qv: QueryVars{Order: []string{"-Create", "-Updated"}}, act: query.Order("-Create").Order("-Updated")},
		{qv: QueryVars{Limit: 10, Offset: 0, Order: []string{"-Create"}}, act: query.Limit(10).Offset(0).Order("-Create")},
	}

	for i, test := range tests {
		exp, err := ds.BuildQuery("test", test.qv)
		g.Ok(t, err)
		g.EqualsWithNumber(t, i, exp, test.act)
	}
}

func TestBuildQueryOperator(t *testing.T) {
	t.Parallel()

	g := New()

	ctx := memory.Use(context.Background())
	ds, err := g.NewDatastore(ctx)
	g.Ok(t, err)

	query := ds.NewQuery("test")

	var tests = []struct {
		operator string
		act      *datastore.Query
	}{
		{operator: "=", act: query.Eq("key", "value")},
		{operator: ">", act: query.Lt("key", "value")},
		{operator: ">=", act: query.Lte("key", "value")},
		{operator: "=>", act: query.Lte("key", "value")},
		{operator: "<", act: query.Gt("key", "value")},
		{operator: "<=", act: query.Gte("key", "value")},
		{operator: "=<", act: query.Gte("key", "value")},
	}

	for i, test := range tests {
		filter := Filter{Field: "key", Operator: test.operator, Value: "value"}
		qv := QueryVars{}
		qv.Filter = append(qv.Filter, filter)
		exp, err := ds.BuildQuery("test", qv)
		g.Ok(t, err)
		g.EqualsWithNumber(t, i, exp, test.act)
	}
}
