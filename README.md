# Gochi

Gochi is lightway WAF for Google AppEngine.

## Usage

```
package main

import "github/nobuhito/gochi"

var g *gochi.Gochi

func init() {
  g = gochi.New()
  g.STATIC("public")

  users := g.Group("/api/users")
  g.GET("/:id", getUsers)

  http.Handle("/", g.Router)
}

type User struct {
  ID int `json:"id"`
  Name string `json:"name"`
}

type Users []User

func getUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  res := gochi.NewResponse(w)
  res.Body = Users[p.ByName("id")]
  gochi.LogDebug(r, res.Body)
  res.WriteJSON(w)
}
```
## API

### Router

 - type
   - Group
 - method gochi
   - gochi.Static(path string)
   - gochi.Get(path string, handle httprouter.Handle)
   - gochi.Put(path string, handle httprouter.Handle)
   - gochi.Delete(path string, handle httprouter.Handle)
   - gochi.Post(path string, handle httprouter.Handle)
   - gochi.Group(root string) Group
 - method group
   - group.Get(path string, handle httprouter.Handle)
   - group.Put(path string, handle httprouter.Handle)
   - group.Delete(path string, handle httprouter.Handle)
   - group.Post(path string, handle httprouter.Handle)

### Log

  - method
    - LogDebug(req \*http.Request, message interface{})

### Response

 - type
   - Response
 - method gochi
   - gochi.NewResponse(w http.ResponseWriter) Response
 - method response
   - response.Write(w http.ResponseWriter)
   - response.WriteJSON(w http.ResponseWriter)

### Datastore

 - type
   - Datastore
   - SearchFilter
   - SearchQuery
 - method gochi
   - gochi.NewDatastore(r \*http.Request) Datastore
   - gochi.NewSeachQuery(r \*http.Request) SearchQuery
 - method datastore
   - datastore.Put(data interface{}) error
   - datastore.Del(data interface{}) error
 - method searchQuery
   - searchQuery.GetQuery() (datastore.Query, error)

### Search

 - type
   - FullTextSearchContent
   - FullTextSearch
 - method gochi
   - gochi.NewFullTextSearchContents(id string) FullTextSearchContent
   - gochi.NewFullTextSearch(r \*http.Request, namespace string) FullTextSearch
 - method fullTextSearch
   - fullTextSearch.Put(content FullTextSearchContent) error
   - fullTextSearch.Del(content FullTextSearchContent) error
   - fullTextSearch.Get(query string) (ids []string, err error)
