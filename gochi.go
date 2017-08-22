package gochi

import (
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.chromium.org/gae/impl/memory"
	"go.chromium.org/gae/impl/prod"
	"go.chromium.org/luci/common/logging"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const version = "0.0.1"

type Env int

// Env is environment
const (
	PROD Env = iota
	DEV
	TEST
)

func (env Env) String() string {
	switch env {
	case PROD:
		return "Production"
	case DEV:
		return "Development"
	case TEST:
		return "Test"
	default:
		return "Unknown"
	}
}

type Gochi struct { // Gochi
	Verion string
	InDev  bool
	Router *mux.Router
	//Env     Env
	//Context context.Context
}

func New() *Gochi {
	return &Gochi{
		Verion: version,
		Router: initRouter(),
		//Env:     PROD,
		//Context: nil,
	}
}

func LogDebug(ctx context.Context, message interface{}) {
	log.Debugf(ctx, "%#v", message)
}

func GetContext(appEncineRequest *http.Request) context.Context {
	if os.Getenv("GOCHIENV") == "TEST" || appEncineRequest == nil {
		return memory.Use(context.Background())
	} else {
		return prod.Use(appengine.NewContext(appEncineRequest), appEncineRequest)
	}
}

func ErrorWrap(err error) error {
	return errors.Wrap(err, "")
}

func ErrorWrapWithMessage(err error, message string) error {
	return errors.Wrap(err, message)
}

func ErrorWrapWithMessagef(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func ErrorPrint(ctx context.Context, fmt string, err error) {
	logging.Errorf(ctx, "%+v", errors.WithStack(err))
}
