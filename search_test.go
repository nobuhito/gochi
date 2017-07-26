// +build appengine

package gochi

import (
	"fmt"
	"testing"
	"time"

	"google.golang.org/appengine/search"
)

type indexContent struct {
	ID   string
	HTML search.HTML
	F1   search.Atom
	F2   time.Time
}

func TestSearchWord(t *testing.T) {
	g := New()

	_, ctx, spinDwon := SpinUp(t)
	defer spinDwon()

	namespace := "namespace"
	s := g.NewFullTextSearch(ctx, namespace)

	now := time.Now()
	exps := []indexContent{
		indexContent{
			ID:   "abcdefg",
			HTML: "<div>hogehoge test 日本語</div>",
			F1:   "test1 test2",
			F2:   now,
		},
		indexContent{
			ID:   "1234567",
			HTML: "<div>hogehuga abc English</div>",
			F1:   "test3 test4",
			F2:   now.AddDate(0, 0, -1),
		},
	}

	for _, exp := range exps {
		err := s.Put(exp.ID, &exp)
		g.Ok(t, err)
	}

	var tests = []struct {
		query string
		rows  int
	}{
		{query: "hugahuga", rows: 0},
		{query: "日本語", rows: 1},
		{query: "hogehoge", rows: 1},
		{query: "hogehoge 日本語", rows: 1},
		{query: "test1", rows: 0},

		{query: "F1=\"test1\"", rows: 0},
		{query: "F1=\"test1 test2\"", rows: 1},

		{query: "hogehoge OR hogehuga", rows: 2},
		{query: "hogehoge AND hogehuga", rows: 0},

		{query: fmt.Sprintf("F2=\"%s\"", now.Format("2006-01-02")), rows: 1},
		{query: fmt.Sprintf("F2=\"%s\" hogehoge", now.Format("2006-01-02")), rows: 1},
		{query: fmt.Sprintf("F2=\"%s\"", now.Format("2006/01/02")), rows: 0},
		{query: fmt.Sprintf("F2=\"%s\"", now.Format("2006-01")), rows: 0},

		{
			query: fmt.Sprintf(
				"F2=\"%s\" OR F2=\"%s\"",
				now.Format("2006-01-02"),
				now.AddDate(0, 0, -1).Format("2006-01-02"),
			),
			rows: 2,
		},
	}

	for i, test := range tests {
		act := indexContent{}

		IDs, err := s.Search(test.query, &act)
		g.Ok(t, err)
		g.EqualsWithNumber(t, i, test.rows, len(IDs))
	}

}

func TestSearchCRUD(t *testing.T) {
	g := New()

	_, ctx, spinDwon := SpinUp(t)
	defer spinDwon()

	namespace := "namespace"
	s := g.NewFullTextSearch(ctx, namespace)
	now := time.Now()
	exps := []indexContent{
		indexContent{
			ID:   "abcdefg",
			HTML: "<div>hogehoge test 日本語</div>",
			F2:   now,
		},
	}

	for _, exp := range exps {
		err := s.Put(exp.ID, &exp)
		g.Ok(t, err)
	}

	act := indexContent{}

	err := s.Get(exps[0].ID, &act)
	g.Ok(t, err)
	act.F2 = now
	g.Equals(t, exps[0], act)

	err = s.Del(exps[0].ID)
	g.Ok(t, err)

	act1 := indexContent{}
	_ = s.Get(exps[0].ID, &act1)
	act.F2 = now
	g.Equals(t, indexContent{}, act1)
}
