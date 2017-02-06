package gochi

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/search"
)

type FullTextSearchContent struct {
	ID   string
	Html search.HTML
}

func (g *Gochi) NewFullTextSearchContents(id string) FullTextSearchContent {
	return FullTextSearchContent{
		ID:   id,
		Html: "",
	}
}

type FullTextSearch struct {
	Namespace string
	Request   *http.Request
}

func (g *Gochi) NewFullTextSearch(r *http.Request, namespace string) FullTextSearch {
	return FullTextSearch{
		Namespace: namespace,
		Request:   r,
	}
}

func (s *FullTextSearch) Put(content FullTextSearchContent) error {
	index, err := search.Open(s.Namespace)
	if err != nil {
		return err
	}

	c := appengine.NewContext(s.Request)
	_, err = index.Put(c, content.ID, content)
	if err != nil {
		return err
	}

	return nil
}

func (s *FullTextSearch) Del(content FullTextSearchContent) error {
	index, err := search.Open(s.Namespace)
	if err != nil {
		return err
	}

	c := appengine.NewContext(s.Request)
	return index.Delete(c, content.ID)
}

func (s *FullTextSearch) Get(query string) (ids []string, err error) {
	var results []string
	index, err := search.Open(s.Namespace)
	if err != nil {
		return results, err
	}

	c := appengine.NewContext(s.Request)
	for item := index.Search(c, query, nil); ; {
		var id string
		var content FullTextSearchContent
		id, err = item.Next(content)
		if err == search.Done {
			break
		}
		if err != nil {
			return results, err
		}

		results = append(results, id)
	}

	return results, nil
}
