package gochi

import (
	"context"
	"fmt"
	"os"
	"testing"

	"google.golang.org/appengine/aetest"

	"github.com/favclip/testerator"
)

func TestMain(m *testing.M) {
	_, _, err := testerator.SpinUp()

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	status := m.Run()

	err = testerator.SpinDown()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}

func (g *Gochi) SpinUp(tb testing.TB) (inst aetest.Instance, ctx context.Context, spinDown func()) {
	inst, ctx, err := testerator.SpinUp()
	if err != nil {
		g.Assert(tb, true, "could not spinup testerator.", err)
	}

	return inst, ctx, func() {
		err := testerator.SpinDown()
		if err != nil {
			g.Assert(tb, true, "could not spindown testerator.", err)
		}
	}

}
