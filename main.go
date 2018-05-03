package main

import (
	"github.com/jessevdk/go-flags"
	logptn "github.com/m-mizutani/logptn/lib"
	"log"
	"os"
)

type Options struct {
	// MaxLen uint   `long:"maxlen" description:"Max length of log message"`
	Output    string `short:"o" long:"output" description:"Output file, '-' means stdout" default:"-"`
	OutFormat string `short:"f" long:"format" choice:"text" choice:"json" default:"text"`
	// FileName string `short:"i" description:"A log file" value-name:"FILE"`
}

func main() {
	var opts Options

	args, ParseErr := flags.ParseArgs(&opts, os.Args)
	if ParseErr != nil {
		os.Exit(1)
	}

	var writer logptn.Writer
	switch opts.OutFormat {
	case "text":
		writer = &logptn.TextWriter{}
	case "json":
		writer = &logptn.JsonWriter{}
	}
	writer.Open(opts.Output)

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

	writer.Dump(gen.Formats())
}
