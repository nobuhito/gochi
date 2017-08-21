// +build appengine

package gochi

import (
	"testing"

	"google.golang.org/appengine/search"
)

func TestSearchCRUD(t *testing.T) {
	id := "abcdefc"
	html := "<div>hogehoge test 日本語</div>"
	namespace := "namespace"

	g := New(TEST)
	//
	// _, ctx, spinDwon := SpinUp(t)
	// defer spinDwon()

	a := g.NewFullTextSearchContents(id)
	a.HTML = search.HTML(html)

	s := g.NewFullTextSearch(ctx, namespace)
	err := s.Put(&a)
	g.Ok(t, err)

	b, err := s.Get(id)
	g.Ok(t, err)
	g.Equals(t, a, b)

	ids, err := s.Search("hugahuga")
	g.Ok(t, err)
	g.Equals(t, 0, len(ids))

	ids, err = s.Search("日本語")
	g.Ok(t, err)
	g.Equals(t, 1, len(ids))

	ids, err = s.Search("hogehoge")
	g.Ok(t, err)
	g.Equals(t, 1, len(ids))

	c, err := s.Get(ids[0])
	g.Ok(t, err)
	g.Equals(t, a, c)

	err = s.Del(&c)
	g.Ok(t, err)

	ids, err = s.Search("hogehoge")
	g.Ok(t, err)
	g.Equals(t, 0, len(ids))
}
