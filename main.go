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
	Output       string  `short:"o" long:"output" description:"Output file, '-' means stdout" default:"-"`
	Dumper       string  `short:"d" long:"dumper" choice:"text" choice:"json" choice:"sjson" default:"text"`
	Threshold    float64 `short:"t" long:"threshold" default:"0.7"`
	Delimiters   string  `short:"s" long:"delimiters"`
	Content      string  `short:"c" long:"content" choice:"log" choice:"format" default:"format"`
	DisableRegex bool    `long:"disable-regex"`
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
	case "sjson":
		dumper, dumperErr = dump.NewSimpleJsonDumper(opts.Output)
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
	sp := logptn.NewSimpleSplitter()
	ptn.ReplaceSplitter(sp)

	if opts.Delimiters != "" {
		sp.SetDelim(opts.Delimiters)
	}

	if opts.DisableRegex {
		log.Println("disable regex patterns")
		sp.DisableRegex()
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
	switch opts.Content {
	case "log":
		dumper.DumpLog(ptn.Logs())
	case "format":
		dumper.DumpFormat(ptn.Formats())
	}
}
