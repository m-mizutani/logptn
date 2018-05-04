package main

import (
	"github.com/jessevdk/go-flags"
	logptn "github.com/m-mizutani/logptn/lib"
	"log"
	"os"
)

type options struct {
	// MaxLen uint   `long:"maxlen" description:"Max length of log message"`
	Output     string  `short:"o" long:"output" description:"Output file, '-' means stdout" default:"-"`
	OutFormat  string  `short:"f" long:"format" choice:"text" choice:"json" default:"text"`
	Threshold  float64 `short:"t" long:"threshold" default:"0.7"`
	Delimiters string  `short:"d" long:"delimiters"`
	// FileName string `short:"i" description:"A log file" value-name:"FILE"`
}

func main() {
	var opts options

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

	if werr := writer.Open(opts.Output); werr != nil {
		log.Fatal("File open error: ", werr)
		os.Exit(1)
	}

	gen := logptn.NewGenerator()

	if opts.Delimiters != "" {
		sp := logptn.NewSimpleSplitter()
		sp.SetDelim(opts.Delimiters)
		gen.ReplaceSplitter(sp)
	}

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
