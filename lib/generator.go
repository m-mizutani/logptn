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
	splitter Splitter
	builder  ClusterBuilder
}

// Constructor of Generator
func NewGenerator() *Generator {
	g := Generator{}
	g.splitter = NewSplitter()
	g.builder = NewSimpleClusterBuilder()
	return &g
}

// Replace splitter of Generator
func (x *Generator) ReplaceSplitter(sp Splitter) {
	x.splitter = sp
}

// Replace ClusterBuilder of Generator
func (x *Generator) ReplaceClusterBuilder(builder ClusterBuilder) {
	x.builder = builder
}

// Read lines from log file (not only raw text but also gzip)
func (x *Generator) ReadFile(fpath string) error {
	fp, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer fp.Close()

	var fio io.Reader
	if strings.HasSuffix(fpath, ".gz") {
		zr, zerr := gzip.NewReader(fp)
		if zerr != nil {
			return zerr
		}
		defer zr.Close()
		fio = zr
	} else {
		fio = fp
	}

	return x.ReadIO(fio)
}

// Read lines from io.Reader
func (x *Generator) ReadIO(fp io.Reader) error {
	s := bufio.NewScanner(fp)
	for s.Scan() {
		text := s.Text()
		if len(text) > 0 {
			if err := x.ReadLine(s.Text()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Read a line from argument.
func (x *Generator) ReadLine(msg string) error {
	log := newLog(msg, x.splitter)
	x.logs = append(x.logs, log)
	return nil
}

// Finalize and build format(s).
func (x *Generator) Finalize() {
	clusters := x.builder.Clustering(x.logs)
	for _, cluster := range clusters {
		format := GenFormat(cluster)
		x.formats = append(x.formats, format)
	}
}

// Getter of formats
func (x *Generator) Formats() []*Format {
	return x.formats
}

// Getter of logs
func (x *Generator) Logs() []*Log {
	return x.logs
}
