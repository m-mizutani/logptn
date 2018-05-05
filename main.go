package main

import (
	"github.com/jessevdk/go-flags"
	logptn "github.com/m-mizutani/logptn/lib"
	dump "github.com/m-mizutani/logptn/lib/dumper"
	"log"
	"os"
)

type options struct {
	// MaxLen uint   `long:"maxlen" description:"Max length of log message"`
	Output     string  `short:"o" long:"output" description:"Output file, '-' means stdout" default:"-"`
	Dumper     string  `short:"d" long:"dumper" choice:"text" choice:"json" default:"text"`
	Threshold  float64 `short:"t" long:"threshold" default:"0.7"`
	Delimiters string  `short:"s" long:"delimiters"`

	// FileName string `short:"i" description:"A log file" value-name:"FILE"`
}

func main() {
	var opts options

	args, ParseErr := flags.ParseArgs(&opts, os.Args)
	if ParseErr != nil {
		os.Exit(1)
	}

	// Setup Writer
	var dumper dump.Dumper
	var dumperErr error
	switch opts.Dumper {
	case "text":
		dumper, dumperErr = dump.NewTextDumper(opts.Output)
	case "json":
		dumper, dumperErr = dump.NewJsonDumper(opts.Output)
	default:
		panic("No such dumper: " + opts.Dumper)
	}

	if dumperErr != nil {
		log.Fatal("File open error: ", dumperErr)
		os.Exit(1)
	}

	// Creating pattern generator.
	ptn := logptn.NewPattern()

	// Setup Splitter
	if opts.Delimiters != "" {
		sp := logptn.NewSimpleSplitter()
		sp.SetDelim(opts.Delimiters)
		ptn.ReplaceSplitter(sp)
	}

	// Setup ClusterBuilder
	if opts.Threshold > 0 {
		builder := logptn.NewSimpleClusterBuilder()
		builder.SetThreshold(opts.Threshold)
		ptn.ReplaceClusterBuilder(builder)
	}

	for _, arg := range args[1:] {
		log.Print("arg:", arg)
		var err error

		if arg != "-" {
			err = ptn.ReadFile(arg)
		} else {
			err = ptn.ReadIO(os.Stdin)
		}

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	ptn.Finalize()
	dumper.DumpFormat(ptn.Formats())
}
