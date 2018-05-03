package logptn

import (
// logging "log"
)

type Format struct {
	Chunks []*Chunk
}

func NewFormat(chunks []*Chunk) *Format {
	f := Format{}
	f.Chunks = make([]*Chunk, len(chunks))

	for idx, c := range chunks {
		f.Chunks[idx] = c.Clone()
	}

	return &f
}

func (x *Format) Merge(chunks []*Chunk) {
	for idx, c := range chunks {
		if x.Chunks[idx] != nil && x.Chunks[idx].Data != c.Data {
			x.Chunks[idx] = nil
		}
	}
}

func GenFormat(cluster *Cluster) *Format {
	format := NewFormat(cluster.Base.Chunk)

	for _, log := range cluster.logs {
		format.Merge(log.Chunk)
	}

	return format
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
