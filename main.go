package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	logptn "github.com/m-mizutani/logptn/lib"
	"log"
	"os"
)

type Options struct {
	// MaxLen uint   `long:"maxlen" description:"Max length of log message"`
	Output string `short:"o" long:"output" description:"Output file, '-' means stdout" default:"-"`
	// FileName string `short:"i" description:"A log file" value-name:"FILE"`
}

func main() {
	var opts Options

	args, ParseErr := flags.ParseArgs(&opts, os.Args)
	if ParseErr != nil {
		os.Exit(1)
	}

	var out *os.File
	if opts.Output == "-" {
		out = os.Stdout
	} else {
		fd, err := os.OpenFile(opts.Output, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal("Can not open file: ", opts.Output)
			os.Exit(1)
		}
		defer fd.Close()
		out = fd
	}

	gen := logptn.NewGenerator()

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

	gen.Finalize()

	for _, format := range gen.Formats() {
		fmt.Fprint(out, format.String())
	}
}
