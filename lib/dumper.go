package logptn

import (
	"log"
	"os"
)

// Dump logs and format
type Dumper interface {
	Setup(arg string) error
	DumpFormat(formats []*Format) error
	DumpLog(logs []*Log, formats []*Format) error
	Teardown() error
}

type fileDumper struct {
	out         *os.File
	shouldClose bool
}

func (x *fileDumper) Setup(arg string) error {
	x.shouldClose = false
	if arg == "-" {
		x.out = os.Stdout
	} else {
		fd, err := os.OpenFile(arg, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal("Can not open file: ", arg)
			return err
		}
		x.out = fd
		x.shouldClose = true
	}

	return nil
}

func (x *fileDumper) Teardown(arg string) error {
	if x.shouldClose {
		return x.out.Close()
	}

	return nil
}
