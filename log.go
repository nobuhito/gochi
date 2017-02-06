package gochi

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func LogDebug(req *http.Request, message interface{}) {
	ctx := appengine.NewContext(req)
	log.Debugf(ctx, "%#v", message)
}
