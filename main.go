package main

import (
	"github.com/jessevdk/go-flags"
	logptn "github.com/m-mizutani/logptn/lib"
	"log"
	"os"
)

type Options struct {
	MaxLen uint `long:"maxlen" description:"Max length of log message"`
	// FileName string `short:"i" description:"A log file" value-name:"FILE"`
}

func main() {
	var opts Options

	args, ParseErr := flags.ParseArgs(&opts, os.Args)
	if ParseErr != nil {
		os.Exit(1)
	}

	gen := logptn.Generator{}

	for _, arg := range args[1:] {
		log.Print("arg:", arg)
		var err error

		if arg != "-" {
			err = gen.ReadFile(arg)
		} else {
			err = gen.ReadIO(os.Stdin)
		}

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}
