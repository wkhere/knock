package main

import (
	"fmt"
	"io"
	"os/exec"
)

func parseArgs(args []string) (c config, err error) {

	rest := make([]string, 0, len(args))
flags:
	for ; len(args) > 0; args = args[1:] {
		switch arg := args[0]; {

		case arg == "-v", arg == "--verbose":
			c.verbose = true

		case arg == "-s", arg == "--strict":
			c.strict = true

		case arg == "-h", arg == "--help":
			c.help = func(w io.Writer) { fmt.Fprintln(w, usage) }
			return c, nil

		case arg == "--":
			rest = append(rest, args[1:]...)
			break flags

		case len(arg) > 1 && arg[0] == '-':
			err = fmt.Errorf("unknown flag %s", arg)
			// note: if there will be more possible flag errors,
			// there needs to be local errorf func saving only the 1st err

		default:
			rest = append(rest, arg)
		}
	}

	if err != nil {
		return c, err
	}
	if len(rest) == 0 {
		return c, fmt.Errorf("expecting program to run")
	}
	c.path, err = exec.LookPath(rest[0])
	if err != nil {
		return c, err
	}
	c.args = rest[1:]
	// todo: for Go 1.19 and above, check ErrDot
	// (see https://pkg.go.dev/os/exec@go1.19.5)

	return c, nil
}

const usage = `Usage: knock [-v|--verbose] [-s|--strict] program arg...
  --strict: extra files in the same dir are ignored`
