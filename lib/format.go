package logptn

import (
	// logging "log"
	"fmt"
	"github.com/fatih/color"
	"github.com/satori/go.uuid"
)

// Segment is a part of format
type Segment interface {
	Text() string
	Fixed() bool
	Count() int
	merge(s string) bool
}

// fixed part
type fixture struct {
	Data  string `json:"text"`
	total int
}

func (x *fixture) Text() string {
	return x.Data
}

func (x *fixture) Fixed() bool {
	return true
}

func (x *fixture) Count() int {
	return x.total
}

func (x *fixture) merge(s string) bool {
	if x.Data != s {
		return false
	}

	x.total++
	return true
}

func newFixture(text string) *fixture {
	return &fixture{Data: text}
}

// variable part
type variable struct {
	Values map[string]int `json:"values"`
}

func (x *variable) Text() string {
	return "*"
}

func (x *variable) Fixed() bool {
	return false
}

func (x *variable) Count() int {
	total := 0
	for _, c := range x.Values {
		total += c
	}
	return total
}

func (x *variable) merge(s string) bool {
	x.Values[s]++
	return true
}

func newVariable(f Segment) *variable {
	v := &variable{Values: map[string]int{}}
	v.Values[f.Text()] = f.Count()
	return v
}

// Format is a structure of log format.
type Format struct {
	uuid     string
	Segments []Segment `json:"segments"`
	Count    int       `json:"count"`
	Sample   string    `json:"sample"`
}

func (x Format) String() string {
	red := color.New(color.FgRed).SprintFunc()

	str := fmt.Sprintf("%6d [%s] ", x.Count, x.id())

	for _, s := range x.Segments {
		if s.Fixed() {
			str += s.Text()
		} else {
			str += red(s.Text())
		}
	}

	return str
}

// GenFormat generates a format from Cluster (set of logs)
func GenFormat(cluster Cluster) *Format {
	format := newFormat(cluster.Logs()[0].Chunk)

	for _, log := range cluster.Logs() {
		format.merge(log.Chunk)
		log.format = format
	}

	format.Sample = format.String()

	return format
}

func newFormat(chunks []*Chunk) *Format {
	f := Format{}
	f.Segments = make([]Segment, len(chunks))
	f.uuid = uuid.NewV4().String()

	for idx, c := range chunks {
		f.Segments[idx] = newFixture(c.Data)
	}

	return &f
}

func (x *Format) id() string {
	return x.uuid[:8]
}

func (x *Format) merge(chunks []*Chunk) {
	x.Count++
	for idx, c := range chunks {
		if !x.Segments[idx].merge(c.Data) {
			x.Segments[idx] = newVariable(x.Segments[idx])
			x.Segments[idx].merge(c.Data)
		}
	}
}
