package gochi

import (
	"go.chromium.org/gae/service/info"
	"golang.org/x/net/context"

	"google.golang.org/appengine/search"
)

type FullTextSearchContent struct {
	ID   string
	HTML search.HTML
}

func (g *Gochi) NewFullTextSearchContents(id string) FullTextSearchContent {
	return FullTextSearchContent{
		ID:   id,
		HTML: "",
	}
}

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

func (s *FullTextSearch) Put(content *FullTextSearchContent) error {
	if info := info.GetTestable(s.Context); info != nil { // テスト時は機能させない
		return nil
	}

	index, err := search.Open(s.Namespace)
	if err != nil {
		return err
	}

	//c := appengine.NewContext(s.Request)
	_, err = index.Put(s.Context, content.ID, content)
	if err != nil {
		return err
	}

	return nil
}

func (s *FullTextSearch) Get(id string) (FullTextSearchContent, error) {
	if info := info.GetTestable(s.Context); info != nil { // テスト時は機能させない
		return FullTextSearchContent{}, nil
	}

	var content FullTextSearchContent

	index, err := search.Open(s.Namespace)
	if err != nil {
		return content, ErrorWrap(err)
	}

	// c := appengine.NewContext(s.Request)
	err = index.Get(s.Context, id, &content)
	if err != nil {
		return content, ErrorWrap(err)
	}

	return content, nil
}

func (s *FullTextSearch) Search(query string) (ids []string, err error) {
	if info := info.GetTestable(s.Context); info != nil { // テスト時は機能させない
		return []string{}, nil
	}

	var results []string
	index, err := search.Open(s.Namespace)
	if err != nil {
		return results, ErrorWrap(err)
	}

	// c := appengine.NewContext(s.Request)
	for item := index.Search(s.Context, query, nil); ; {
		var id string
		var content FullTextSearchContent
		id, err = item.Next(&content)
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

func (s *FullTextSearch) Del(content *FullTextSearchContent) error {
	if info := info.GetTestable(s.Context); info != nil { // テスト時は機能させない
		return nil
	}

	index, err := search.Open(s.Namespace)
	if err != nil {
		return err
	}

	// c := appengine.NewContext(s.Request)
	return index.Delete(s.Context, content.ID)
}
