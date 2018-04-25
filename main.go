package main

import (
	"github.com/jessevdk/go-flags"
	genlogfmt "github.com/m-mizutani/genlogfmt/lib"
	"log"
	"os"
)

type Options struct {
	FileName string `short:"r" description:"A log file" value-name:"FILE"`
}

func main() {
	var opts Options

	_, ParseErr := flags.ParseArgs(&opts, os.Args)
	if ParseErr != nil {
		os.Exit(1)
	}

	gen := genlogfmt.Generator{}

	if opts.FileName != "" {
		err := gen.ReadFile(opts.FileName)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}
