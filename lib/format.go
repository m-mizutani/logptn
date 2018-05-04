package logptn

import (
// logging "log"
)

// Segment is a part of format
type Segment interface {
	text() string
	fixed() bool
	count() int
	merge(s string) bool
}

// fixed part
type fixture struct {
	Text  string `json:"text"`
	total int
}

func (x *fixture) text() string {
	return x.Text
}

func (x *fixture) fixed() bool {
	return true
}

func (x *fixture) count() int {
	return x.total
}

func (x *fixture) merge(s string) bool {
	if x.Text != s {
		return false
	}

	x.total++
	return true
}

func newFixture(text string) *fixture {
	return &fixture{Text: text}
}

// variable part
type variable struct {
	Values map[string]int `json:"values"`
}

func (x *variable) text() string {
	return "*"
}

func (x *variable) fixed() bool {
	return false
}

func (x *variable) count() int {
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
	v.Values[f.text()] = f.count()
	return v
}

// Format is a structure of log format.
type Format struct {
	Segments []Segment `json:"segments"`
	Count    int       `json:"count"`
	Sample   string    `json:"sample"`
}

func (x Format) String() string {
	var str string
	for _, s := range x.Segments {
		str += s.text()
	}

	return str
}

// GenFormat generates a format from Cluster (set of logs)
func GenFormat(cluster Cluster) *Format {
	format := newFormat(cluster.Logs()[0].Chunk)

	for _, log := range cluster.Logs() {
		format.merge(log.Chunk)
	}

	format.Sample = format.String()

	return format
}

func newFormat(chunks []*Chunk) *Format {
	f := Format{}
	f.Segments = make([]Segment, len(chunks))

	for idx, c := range chunks {
		f.Segments[idx] = newFixture(c.Data)
	}

	return &f
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
