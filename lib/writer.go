package logptn

import (
	"fmt"
	"log"
	"os"
)

type Writer interface {
	Open(fpath string) error
	Dump(formats []*Format) error
}

type FileWriter struct {
	out *os.File
}

func (x *FileWriter) Open(fpath string) error {
	if fpath == "-" {
		x.out = os.Stdout
	} else {
		fd, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal("Can not open file: ", fpath)
			return err
		}
		x.out = fd
	}

	return nil
}

type TextWriter struct {
	FileWriter
}

func (x *TextWriter) Dump(formats []*Format) error {
	for _, format := range formats {
		fmt.Fprint(x.out, format.String(), "\n")
	}
	return nil
}
