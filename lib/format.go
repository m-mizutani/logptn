package logptn

import (
// logging "log"
)

// Segment is a part of format
type Segment struct {
	Text   string         `json:"text"`
	Values map[string]int `json:"values"`
	Fixed  bool
}

func newSegment(text string) *Segment {
	s := Segment{Text: text, Values: map[string]int{}, Fixed: true}
	return &s
}

// Format is a structure of log format.
type Format struct {
	Segments []*Segment `json:"segments"`
	Count    int
}

func (x Format) String() string {
	var str string
	for _, s := range x.Segments {
		if s.Fixed {
			str += s.Text
		} else {
			str += "*"
		}
	}

	return str
}

// GenFormat generates a format from Cluster (set of logs)
func GenFormat(cluster *Cluster) *Format {
	format := newFormat(cluster.Base.Chunk)

	for _, log := range cluster.logs {
		format.merge(log.Chunk)
	}

	return format
}

func newFormat(chunks []*Chunk) *Format {
	f := Format{}
	f.Segments = make([]*Segment, len(chunks))

	for idx, c := range chunks {
		f.Segments[idx] = newSegment(c.Data)
	}

	return &f
}

func (x *Format) merge(chunks []*Chunk) {
	x.Count++
	for idx, c := range chunks {
		if x.Segments[idx].Fixed && x.Segments[idx].Text != c.Data {
			x.Segments[idx].Fixed = false
		}

		x.Segments[idx].Values[c.Data]++
	}
}
