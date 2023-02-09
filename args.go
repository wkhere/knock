package main

import (
	"fmt"
	"os/exec"
)

func parseArgs(args []string) (c config, _ error) {
	if len(args) < 1 {
		return c, fmt.Errorf(usage)
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		return c, err
	}
	// todo: for Go 1.19 and above, check ErrDot
	// (see https://pkg.go.dev/os/exec@go1.19.5)
	c.path = path
	c.args = args[1:]

	return c, nil
}

const usage = `Usage: knock program arg...`
