package main

import (
	"log"
	"os"
)

type config struct {
	path string
	args []string
}

func main() {
	log.SetFlags(0)

	c, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Print(err)
		os.Exit(2)
	}

	log.SetFlags(log.Ldate | log.Ltime)

	err = run(&c)
	if err != nil {
		log.Fatal(err)
	}
}
