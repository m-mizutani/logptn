package logptn

import (
// logging "log"
)

// Format is a structure of log format.
type Format struct {
	Chunks []*Chunk
	values []map[string]int
	Count  int
}

func (x Format) String() string {
	var s string
	for _, c := range x.Chunks {
		if c != nil {
			s += c.Data
		} else {
			s += "*"
		}
	}

	return s
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
	f.Chunks = make([]*Chunk, len(chunks))
	f.values = make([]map[string]int, len(chunks))

	for idx, c := range chunks {
		f.Chunks[idx] = c.Clone()
		f.values[idx] = map[string]int{}
	}

	return &f
}

func (x *Format) merge(chunks []*Chunk) {
	x.Count++
	for idx, c := range chunks {
		if x.Chunks[idx] != nil && x.Chunks[idx].Data != c.Data {
			x.Chunks[idx] = nil
		}

		x.values[idx][c.Data]++
	}
}
