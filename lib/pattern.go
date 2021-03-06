package logptn

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

type Pattern struct {
	logs     []*Log
	formats  []*Format
	splitter Splitter
	builder  ClusterBuilder
	index    int
}

// Constructor of Pattern
func NewPattern() *Pattern {
	g := Pattern{}
	g.splitter = NewSplitter()
	g.builder = NewSimpleClusterBuilder()
	return &g
}

// Replace splitter of Pattern
func (x *Pattern) ReplaceSplitter(sp Splitter) {
	x.splitter = sp
}

// Replace ClusterBuilder of Pattern
func (x *Pattern) ReplaceClusterBuilder(builder ClusterBuilder) {
	x.builder = builder
}

// Read lines from log file (not only raw text but also gzip)
func (x *Pattern) ReadFile(fpath string) error {
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
func (x *Pattern) ReadIO(fp io.Reader) error {
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
func (x *Pattern) ReadLine(msg string) error {
	log := newLog(msg, x.splitter, x.index)
	x.index++
	x.logs = append(x.logs, log)
	return nil
}

// Finalize and build format(s).
func (x *Pattern) Finalize() {
	clusters := x.builder.Clustering(x.logs)
	for _, cluster := range clusters {
		format := GenFormat(cluster)
		x.formats = append(x.formats, format)
	}
}

// Getter of formats
func (x *Pattern) Formats() []*Format {
	return x.formats
}

// Getter of logs
func (x *Pattern) Logs() []*Log {
	return x.logs
}
