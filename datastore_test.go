package gochi

import (
	"strconv"
	"testing"
)

type Hogehoge struct {
	ID   string `goon:"id"`
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
	g := New()

	_, ctx, spinDwon := g.SpinUp(t)
	defer spinDwon()

	ds := g.NewDatastore(ctx)

	hoge := hoge()

	err := ds.Put(&hoge)
	g.Ok(t, err)

	hogehoge := Hogehoge{
		ID: "abcdefg",
	}

	err = ds.Get(&hogehoge)
	g.Ok(t, err)
	g.Equals(t, hoge, hogehoge)

	err = ds.Del(&hoge)
	g.Ok(t, err)

	e := "datastore: no such entity"
	err = ds.Get(&hoge)
	g.Equals(t, e, err.Error())
}

func TestGetNamespace(t *testing.T) {
	g := New()

	_, ctx, spinDwon := g.SpinUp(t)
	defer spinDwon()

	ds := g.NewDatastore(ctx)

	a := "Hogehoge"
	b := ds.GetNamespace(hoge())
	g.Equals(t, a, b)
}

func TestSearchQuerySet(t *testing.T) {
	g := New()

	_, ctx, spinDwon := g.SpinUp(t)
	defer spinDwon()

	ds := g.NewDatastore(ctx)
	hoge := hoge()

	sq := g.NewSeachQuery()

	namespace := ds.GetNamespace(hoge)
	sq.Kind = namespace
	g.Equals(t, namespace, sq.Kind)

	order := "ID"
	var orders []string
	orders = append(orders, order)
	sq.SetOrder(order)
	g.Equals(t, orders, sq.Orders)

	limit := 10
	sq.SetLimit(limit)
	g.Equals(t, limit, sq.Limit)

	filter := SearchFilter{
		Key:   "ID =",
		Value: "abcdefg",
	}
	var filters []SearchFilter
	filters = append(filters, filter)
	sq.AppendFilter(filter)
	g.Equals(t, filters, sq.Filters)

	err := ds.Put(&hoge)
	g.Ok(t, err)

	q, err := sq.GetQuery()
	g.Ok(t, err)

	var data Hogehoges
	_, err = ds.Goon.GetAll(&q, &data)
	g.Ok(t, err)
	g.Equals(t, 1, len(data))
}

func TestPutMulti(t *testing.T) {
	g := New()

	_, ctx, spindown := g.SpinUp(t)
	defer spindown()

	ds := g.NewDatastore(ctx)

	var hogehoges []Hogehoge
	count := 1000
	for i := 0; i < count; i++ {
		hoge := hoge()
		hoge.ID = "id" + strconv.Itoa(i)
		hogehoges = append(hogehoges, hoge)
	}
	keys, err := ds.PutMulti(&hogehoges)
	g.Ok(t, err)
	g.Equals(t, count, len(keys))
}

func TestGetBindAll(t *testing.T) {
	g := New()

	_, ctx, spindown := g.SpinUp(t)
	defer spindown()

	ds := g.NewDatastore(ctx)

	var hogehoges []Hogehoge
	hogehoges = append(hogehoges, Hogehoge{ID: "a", Abcd: "hogehoge", Efg: 10})
	hogehoges = append(hogehoges, Hogehoge{ID: "b", Abcd: "hogehoge", Efg: 0})
	hogehoges = append(hogehoges, Hogehoge{ID: "c", Abcd: "", Efg: 10})
	hogehoges = append(hogehoges, Hogehoge{ID: "d", Abcd: "", Efg: 0})
	hogehoges = append(hogehoges, Hogehoge{ID: "e", Abcd: "hugahuga", Efg: 1})
	keys, err := ds.PutMulti(&hogehoges)
	g.Ok(t, err)
	g.Equals(t, len(hogehoges), len(keys))

	sq := g.NewSeachQuery()
	namespace := ds.GetNamespace(Hogehoge{})
	sq.Kind = namespace

	var filters []SearchFilter
	filters = append(filters, SearchFilter{Key: "Abcd =", Value: "hogehoge"})
	filters = append(filters, SearchFilter{Key: "Efg =", Value: 10})

	keys, err = ds.GetBindAll(sq, filters, &hogehoges)
	g.Ok(t, err)
	g.Equals(t, 3, len(keys))

}

func TestDelAll(t *testing.T) {
	g := New()

	_, ctx, spindown := g.SpinUp(t)
	defer spindown()

	ds := g.NewDatastore(ctx)

	var hoge Hogehoge
	namespace := ds.GetNamespace(hoge)
	err := ds.DelAll(namespace)
	g.Ok(t, err)

	sq := g.NewSeachQuery()
	sq.SetKind(namespace)
	q, err := sq.GetQuery()
	g.Ok(t, err)

	var hogehoges Hogehoges
	q.KeysOnly()
	keys, err := ds.Goon.GetAll(&q, &hogehoges)
	g.Equals(t, 0, len(keys))

}
