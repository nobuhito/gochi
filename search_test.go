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

func getTestData(now time.Time) []indexContent {
	return []indexContent{
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
		indexContent{
			ID:   "234567890",
			HTML: "<div>hogehoge def English</div>",
			F2:   now.AddDate(0, 0, -2),
		},
	}
}

var namespace = "namespace"

func TestSearchWithOptions(t *testing.T) {
	g := New()
	_, ctx, spindown := SpinUp(t)
	defer spindown()

	s := g.NewFullTextSearch(ctx, namespace)
	now := time.Now()
	exps := getTestData(now)
	for _, exp := range exps {
		err := s.Put(exp.ID, &exp)
		g.Ok(t, err)
	}

	query := "hogehoge"
	ids, _, err := s.Search(query, nil)
	g.Ok(t, err)
	g.Equals(t, 2, len(ids))

	options := search.SearchOptions{
		Limit:   1,
		IDsOnly: true,
		Sort: &search.SortOptions{
			Expressions: []search.SortExpression{
				{Expr: "F2", Reverse: false},
			},
		},
	}
	ids, cursor, err := s.SearchWithOptions(query, nil, &options)
	g.Ok(t, err)
	g.Equals(t, 1, len(ids))
	g.Equals(t, "abcdefg", ids[0])
	fmt.Println(cursor)

	options.Cursor = cursor
	ids, cursor, err = s.SearchWithOptions(query, nil, &options)
	g.Ok(t, err)
	g.Equals(t, 1, len(ids))
	g.Equals(t, "234567890", ids[0])

	options.Cursor = cursor
	ids, cursor, err = s.SearchWithOptions(query, nil, &options)
	g.Ok(t, err)
	g.Equals(t, 0, len(ids))
}

func TestSearchWord(t *testing.T) {
	g := New()

	_, ctx, spinDwon := SpinUp(t)
	defer spinDwon()

	s := g.NewFullTextSearch(ctx, namespace)

	now := time.Now()

	exps := getTestData(now)

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
		{query: "hogehoge", rows: 2},
		{query: "hogehoge 日本語", rows: 1},
		{query: "test1", rows: 0},

		{query: "F1=\"test1\"", rows: 0},
		{query: "F1=\"test1 test2\"", rows: 1},

		{query: "hogehoge OR hogehuga", rows: 3},
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

		IDs, _, err := s.Search(test.query, &act)
		g.Ok(t, err)
		g.EqualsWithNumber(t, i, test.rows, len(IDs))
	}

}

func TestSearchCRUD(t *testing.T) {
	g := New()

	_, ctx, spinDwon := SpinUp(t)
	defer spinDwon()

	s := g.NewFullTextSearch(ctx, namespace)
	now := time.Now()
	exps := getTestData(now)

	for _, exp := range exps {
		err := s.Put(exp.ID, &exp)
		g.Ok(t, err)
	}

	act := indexContent{}
	exp := exps[0]

	err := s.Get(exp.ID, &act)
	g.Ok(t, err)
	act.F2 = now
	g.Equals(t, exp, act)

	err = s.Del(exp.ID)
	g.Ok(t, err)

	act1 := indexContent{}
	_ = s.Get(exp.ID, &act1)
	act.F2 = now
	g.Equals(t, indexContent{}, act1)
}
