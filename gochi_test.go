package gochi

import (
	"fmt"
	"os"
	"testing"

	"github.com/favclip/testerator"
)

func TestMain(m *testing.M) {
	env := "APPENGINE_DEV_APPSERVER"

	inDevServer := true
	if os.Getenv(env) == "" {
		inDevServer = false
	}

	if !inDevServer {
		// fmt.Printf("%s not set\n", env)
		// os.Exit(1)
		status := m.Run()
		os.Exit(status)

	} else {
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
}
