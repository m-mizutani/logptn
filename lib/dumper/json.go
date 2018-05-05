package dumper

import (
	"encoding/json"
	"fmt"
	logptn "github.com/m-mizutani/logptn/lib"
)

type jsonOutputTemplate struct {
	Formats []*logptn.Format `json:"formats"`
}

// Constructor of JsonDumper
func NewJsonDumper(fpath string) (*JsonDumper, error) {
	dumper := JsonDumper{}
	if err := dumper.open(fpath); err != nil {
		return nil, err
	}
	return &dumper, nil
}

type JsonDumper struct {
	fileDumper
}

// Dump format data as json structured data
func (x *JsonDumper) DumpFormat(formats []*logptn.Format) error {
	d := jsonOutputTemplate{}
	d.Formats = formats
	b, _ := json.Marshal(d)
	fmt.Fprint(x.out, string(b))
	return nil
}

// Dump format data as json structured data
func (x *JsonDumper) DumpLog(logs []*logptn.Log) error {
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}
