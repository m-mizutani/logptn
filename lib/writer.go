package logptn

import (
	"encoding/json"
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

// Dump format data as plain text
func (x *TextWriter) Dump(formats []*Format) error {
	for _, format := range formats {
		fmt.Fprint(x.out, format.String(), "\n")
	}
	return nil
}

type JsonWriter struct {
	FileWriter
}

type jsonOutputTemplate struct {
	Formats []*Format `json:"formats"`
}

// Dump format data as json structured data
func (x *JsonWriter) Dump(formats []*Format) error {
	d := jsonOutputTemplate{}
	d.Formats = formats
	b, _ := json.Marshal(d)
	fmt.Fprint(x.out, string(b))
	return nil
}
