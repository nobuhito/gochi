package gochi

import (
	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"google.golang.org/appengine/log"
)

const version = "0.0.1"

type Gochi struct { // Gochi
	Verion  string
	InDev   bool
	Router  *mux.Router
	Tasks   []task
	Context context.Context
}

func New() *Gochi {
	return &Gochi{
		Verion:  version,
		Router:  initRouter(),
		InDev:   false,
		Tasks:   []task{},
		Context: nil,
	}
}

func LogDebug(ctx context.Context, message interface{}) {
	log.Debugf(ctx, "%#v", message)
}
