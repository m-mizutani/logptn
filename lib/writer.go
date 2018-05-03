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

func (x *TextWriter) Dump(formats []*Format) error {
	for _, format := range formats {
		fmt.Fprint(x.out, format.String(), "\n")
	}
	return nil
}

type JsonWriter struct {
	FileWriter
}

type JsonOutput struct {
	Formats []*PrintableFormat `json:"formats"`
}

type PrintableFormat struct {
	Literals []*string `json:"literals"`
}

func (x *JsonWriter) Dump(formats []*Format) error {
	d := JsonOutput{}
	for _, format := range formats {
		pf := &PrintableFormat{}
		var s, p *string
		for _, c := range format.Chunks {
			if c != nil {
				s = &c.Data
			} else {
				s = nil
			}

			if p != nil && s != nil {
				t := *p + *s
				pf.Literals[len(pf.Literals)-1] = &t
				p = &t
			} else {
				pf.Literals = append(pf.Literals, s)
				p = s
			}
		}

		d.Formats = append(d.Formats, pf)
	}

	b, _ := json.Marshal(d)
	fmt.Fprint(x.out, string(b))
	return nil
}
