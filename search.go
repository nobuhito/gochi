package gochi

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/search"
)

// type FullTextSearchContent struct {
// 	ID   string
// 	HTML search.HTML
// }

// func (g *Gochi) NewFullTextSearchContents(id string) FullTextSearchContent {
// 	return FullTextSearchContent{
// 		ID:   id,
// 		HTML: "",
// 	}
// }

type FullTextSearch struct {
	Namespace string
	Context   context.Context
}

func (g *Gochi) NewFullTextSearch(ctx context.Context, namespace string) FullTextSearch {
	return FullTextSearch{
		Namespace: namespace,
		Context:   ctx,
	}
}

func (s *FullTextSearch) Put(id string, dst interface{}) error {
	index, err := search.Open(s.Namespace)
	if err != nil {
		return err
	}

	id, err = index.Put(s.Context, id, dst)
	if err != nil {
		return err
	}

	return nil
}

func (s *FullTextSearch) Get(id string, dst interface{}) error {

	index, err := search.Open(s.Namespace)
	if err != nil {
		return ErrorWrap(err)
	}

	err = index.Get(s.Context, id, dst)
	if err != nil {
		return ErrorWrap(err)
	}

	return nil
}

func (s *FullTextSearch) Search(query string, dst interface{}) (ids []string, err error) {
	return s.SearchWithOptions(query, dst, nil)
}

func (s *FullTextSearch) SearchWithOptions(query string, dst interface{}, options *search.SearchOptions) (ids []string, err error) {

	var results []string
	index, err := search.Open(s.Namespace)
	if err != nil {
		return results, ErrorWrap(err)
	}

	for item := index.Search(s.Context, query, options); ; {
		var id string
		id, err = item.Next(dst)
		if err == search.Done {
			break
		}
		if err != nil {
			return results, ErrorWrap(err)
		}

		results = append(results, id)
	}

	return results, nil
}

func (s *FullTextSearch) Del(id string) error {
	index, err := search.Open(s.Namespace)
	if err != nil {
		return err
	}

	return index.Delete(s.Context, id)
}
