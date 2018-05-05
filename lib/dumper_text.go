package logptn

import (
	"fmt"
)

type TextDumper struct {
	fileDumper
}

func (x *TextDumper) DumpFormat(formats []*Format) error {
	for _, format := range formats {
		fmt.Println(format)
	}
	return nil
}

func (x *TextDumper) DumpLog(logs []*Log) error {
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}
