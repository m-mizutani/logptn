package logptn

import (
	"bufio"
	"log"
	"os"
)

type Generator struct {
	logs    []*Log
	formats []*Format
}

func (x *Generator) ReadFile(fpath string) error {
	fp, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer fp.Close()

	return x.ReadIO(fp)
}

func (x *Generator) ReadIO(fp *os.File) error {
	log.Print("ReadIO")
	s := bufio.NewScanner(fp)
	for s.Scan() {
		text := s.Text()
		if len(text) > 0 {
			x.ReadLine(s.Text())
		}
	}
	return nil
}

func (x *Generator) ReadLine(msg string) error {
	log := NewLog(msg)
	x.logs = append(x.logs, log)
	return nil
}

func (x *Generator) Finalize() {
	clusters := Clustering(x.logs)
	for _, cluster := range clusters {
		format := GenFormat(&cluster)
		x.formats = append(x.formats, format)
	}
}

func (x *Generator) Formats() []*Format {
	return x.formats
}

func (x *Generator) Logs() []*Log {
	return x.logs
}
