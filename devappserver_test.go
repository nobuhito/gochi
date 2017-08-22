package gochi

import (
	"context"
	"testing"

	"github.com/favclip/testerator"
	"google.golang.org/appengine/aetest"
)

func TestDevAppServer(t *testing.T) {
}

func SpinUp(tb testing.TB) (inst aetest.Instance, ctx context.Context, spinDown func()) {
	g := New()
	// g.SetTestContext()

	inst, ctx, err := testerator.SpinUp()
	if err != nil {
		g.Assert(tb, true, "could not SpinUp testerator.", err)
	}

	return inst, ctx, func() {
		err := testerator.SpinDown()
		if err != nil {
			g.Assert(tb, true, "could not spindown testerator.", err)
		}
	}
}
