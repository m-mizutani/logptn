package dumper

import (
	logptn "github.com/m-mizutani/logptn/lib"
	"log"
	"os"
)

// Dump logs and format
type Dumper interface {
	DumpFormat(formats []*logptn.Format) error
	DumpLog(logs []*logptn.Log) error
	Shutdown() error
}

type fileDumper struct {
	out         *os.File
	shouldClose bool
}

func (x *fileDumper) open(arg string) error {
	x.shouldClose = false
	if arg == "-" {
		x.out = os.Stdout
	} else {
		fd, err := os.Create(arg)
		if err != nil {
			log.Fatal("Can not open file: ", arg)
			return err
		}
		x.out = fd
		x.shouldClose = true
	}

	return nil
}

func (x *fileDumper) close() error {
	if x.shouldClose {
		return x.out.Close()
	}

	return nil
}

func (x *fileDumper) Shutdown() error {
	return x.close()
}
