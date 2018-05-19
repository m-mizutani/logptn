package logptn

import (
	"fmt"
	"github.com/fatih/color"
)

type Chunk struct {
	Data   string
	Freeze bool
}

type Log struct {
	Index  int
	text   string
	Chunk  []*Chunk
	format *Format
}

func newLog(line string, sp Splitter, idx int) *Log {
	log := Log{}
	log.text = line
	log.Chunk = sp.Split(line)
	log.Index = idx
	return &log
}

func (x *Log) String() string {
	red := color.New(color.FgRed).SprintFunc()

	s := fmt.Sprintf("[%s] ", x.format.id())
	for idx, c := range x.Chunk {
		if x.format.Segments[idx].Fixed() {
			s += c.Data
		} else {
			s += red(c.Data)
		}
	}
	return s
}

func newChunk(d string) *Chunk {
	c := Chunk{}
	c.Data = d
	c.Freeze = false
	return &c
}

func (x *Chunk) Clone() *Chunk {
	c := newChunk(x.Data)
	c.Freeze = x.Freeze
	return c
}

func (x *Chunk) String() string {
	return x.Data
}
