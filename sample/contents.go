package main

func contents() {
	search := g.Group("/api/Contents")
	search.Get("/:query", searchFile)
}
