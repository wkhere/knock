package main

import (
	"io"
	"log"
	"os"
)

type config struct {
	path string
	args []string

	strict  bool
	compare bool
	verbose bool

	help func(io.Writer)
}

func main() {
	log.SetFlags(0)

	c, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Print(err)
		os.Exit(2)
	}
	if c.help != nil {
		c.help(os.Stdout)
		os.Exit(0)
	}

	log.SetFlags(log.Lmicroseconds)

	err = run(&c)
	if err != nil {
		log.Fatal(err)
	}
}
