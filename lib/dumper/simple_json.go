package dumper

import (
	"encoding/json"
	"fmt"
	logptn "github.com/m-mizutani/logptn/lib"
)

type sjFormat struct {
	Segments []*string `json:"segments"`
	Count    int       `json:"count"`
}

type sjOutputTemplate struct {
	Formats []*sjFormat `json:"formats"`
}

// Constructor of JsonDumper
func NewSimpleJsonDumper(fpath string) (*SimpleJsonDumper, error) {
	dumper := SimpleJsonDumper{}
	if err := dumper.open(fpath); err != nil {
		return nil, err
	}
	return &dumper, nil
}

type SimpleJsonDumper struct {
	fileDumper
}

func makeFormatSimple(original *logptn.Format) *sjFormat {
	newFmt := sjFormat{}
	var ptr *string
	for _, seg := range original.Segments {
		if seg.Fixed() {
			if ptr == nil {
				ptr = new(string)
				*ptr = seg.Text()
				newFmt.Segments = append(newFmt.Segments, ptr)
			} else {
				*ptr += seg.Text()
			}
		} else {
			newFmt.Segments = append(newFmt.Segments, nil)
			ptr = nil
		}
	}
	newFmt.Count = original.Count
	return &newFmt
}

// Dump format data as json structured data
func (x *SimpleJsonDumper) DumpFormat(formats []*logptn.Format) error {
	d := sjOutputTemplate{}
	d.Formats = make([]*sjFormat, len(formats))
	for idx, f := range formats {
		d.Formats[idx] = makeFormatSimple(f)
	}

	b, _ := json.Marshal(d)
	fmt.Fprint(x.out, string(b))
	return nil
}

// Dump format data as json structured data
func (x *SimpleJsonDumper) DumpLog(logs []*logptn.Log) error {
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}
