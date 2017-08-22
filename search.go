package gochi

import (
	"context"

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
	//	if info := info.GetTestable(s.Context); info != nil { // テスト時は機能させない
	//		return nil
	//	}

	index, err := search.Open(s.Namespace)
	if err != nil {
		return err
	}

	//c := appengine.NewContext(s.Request)
	_, err = index.Put(s.Context, id, dst)
	if err != nil {
		return err
	}

	return nil
}

func (s *FullTextSearch) Get(id string, dst interface{}) error {
	//	if info := info.GetTestable(s.Context); info != nil { // テスト時は機能させない
	//		return FullTextSearchContent{}, nil
	//	}

	//	var content FullTextSearchContent

	index, err := search.Open(s.Namespace)
	if err != nil {
		return ErrorWrap(err)
	}

	// c := appengine.NewContext(s.Request)
	err = index.Get(s.Context, id, dst)
	if err != nil {
		return ErrorWrap(err)
	}

	return nil
}

func (s *FullTextSearch) Search(query string, dst interface{}) (ids []string, cursor search.Cursor, err error) {
	return s.SearchWithOptions(query, dst, nil)
}

func (s *FullTextSearch) SearchWithOptions(query string, dst interface{}, options *search.SearchOptions) (ids []string, cursor search.Cursor, err error) {
	//	if info := info.GetTestable(s.Context); info != nil { // テスト時は機能させない
	//		return []string{}, nil
	//	}

	// var results []string
	index, err := search.Open(s.Namespace)
	if err != nil {
		return ids, cursor, ErrorWrap(err)
	}

	// c := appengine.NewContext(s.Request)
	for item := index.Search(s.Context, query, options); ; {
		var id string
		id, err = item.Next(dst)
		if err == search.Done {
			break
		}
		if err != nil {
			return ids, cursor, ErrorWrap(err)
		}

		ids = append(ids, id)
		cursor = item.Cursor()
	}

	return ids, cursor, nil
}

func (s *FullTextSearch) Del(id string) error {
	//	if info := info.GetTestable(s.Context); info != nil { // テスト時は機能させない
	//		return nil
	//	}

	index, err := search.Open(s.Namespace)
	if err != nil {
		return err
	}

	// c := appengine.NewContext(s.Request)
	return index.Delete(s.Context, id)
}
