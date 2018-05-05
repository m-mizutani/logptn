package dumper

import (
	"fmt"
	logptn "github.com/m-mizutani/logptn/lib"
)

// Constructor of TextDumper
func NewTextDumper(fpath string) (*TextDumper, error) {
	dumper := TextDumper{}
	if err := dumper.open(fpath); err != nil {
		return nil, err
	}
	return &dumper, nil
}

type TextDumper struct {
	fileDumper
}

func (x *TextDumper) DumpFormat(formats []*logptn.Format) error {
	for _, format := range formats {
		fmt.Println(format)
	}
	return nil
}

func (x *TextDumper) DumpLog(logs []*logptn.Log) error {
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}
