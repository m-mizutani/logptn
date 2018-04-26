package logptn

import (
	// "log"
	"strings"
)

type Splitter struct {
	delims string
}

func NewSplitter() *Splitter {
	s := Splitter{}
	s.delims = " \t!,:;[]{}()=|\\*"
	return &s
}

func (x *Splitter) SetDelim(d string) {
	x.delims = d
}

func (x *Splitter) Split(msg string) []*Chunk {
	var res []*Chunk
	for {
		idx := strings.IndexAny(msg, x.delims)
		if idx < 0 {
			if len(msg) > 0 {
				res = append(res, NewChunk(msg))
			}
			break
		}

		// log.Print("index: ", idx)
		fwd := idx + 1

		s1 := msg[:idx]
		s2 := msg[idx:fwd]
		s3 := msg[fwd:]

		if len(s1) > 0 {
			// log.Print("add s1: ", s1)
			res = append(res, NewChunk(s1))
		}

		if len(s2) > 0 {
			// log.Print("add s2: ", s2)
			res = append(res, NewChunk(s2))
		}

		msg = s3
		// log.Print("remain: ", msg)
	}

	return res
}
