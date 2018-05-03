package logptn

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

type Generator struct {
	logs     []*Log
	formats  []*Format
	splitter *Splitter
}

func NewGenerator() *Generator {
	g := Generator{}
	g.splitter = NewSplitter()
	return &g
}

func (x *Generator) ReadFile(fpath string) error {
	fp, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer fp.Close()

	if strings.HasSuffix(fpath, ".gz") {
		zr, zerr := gzip.NewReader(fp)
		if zerr != nil {
			return zerr
		}
		defer zr.Close()
		return x.ReadIO(zr)
	} else {
		return x.ReadIO(fp)
	}
}

func (x *Generator) ReadIO(fp io.Reader) error {
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
	log := NewLog(msg, x.splitter)
	x.logs = append(x.logs, log)
	return nil
}

func (x *Generator) Finalize() {
	clusters := Clustering(x.logs)
	for _, cluster := range clusters {
		format := GenFormat(cluster)
		x.formats = append(x.formats, format)
	}
}

func (x *Generator) Formats() []*Format {
	return x.formats
}

func (x *Generator) Logs() []*Log {
	return x.logs
}
