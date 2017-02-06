package gochi

import "github.com/julienschmidt/httprouter"

const version = "0.0.1"

type Gochi struct { // Gochi
	Verion string
	InDev  bool
	Router *httprouter.Router
}

func New() *Gochi {
	return &Gochi{
		Verion: version,
		Router: httprouter.New(),
		InDev:  false,
	}
}
