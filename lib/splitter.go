package logptn

import (
	// "log"
	"regexp"
	"strings"
)

type Splitter interface {
	Split(msg string) []*Chunk
}

type SimpleSplitter struct {
	delims    string
	regexList []*regexp.Regexp
	useRegex  bool
}

// NewSplitter is a wrapper of SimpleSplitter
func NewSplitter() Splitter {
	return NewSimpleSplitter()
}

// NewSimpleSplitter is a constructor of SimpleSplitter
func NewSimpleSplitter() *SimpleSplitter {
	s := &SimpleSplitter{}
	s.delims = " \t!,:;[]{}()<>=|\\*\"'"
	s.useRegex = true

	heuristicsPatterns := []string{
		// DateTime
		`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+`,
		`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`,
		`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`,
		// Date
		`\d{4}/\d{2}/\d{2}`,
		`\d{4}-\d{2}-\d{2}`,
		`\d{2}:\d{2}:\d{2}.\d+`,
		// Time
		`\d{2}:\d{2}:\d{2}`,
		// Mail address
		`[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*`,
		// IPv4 address
		`(\d{1,3}\.){3}\d{1,3}`,
	}

	s.regexList = make([]*regexp.Regexp, len(heuristicsPatterns))
	for idx, p := range heuristicsPatterns {
		s.regexList[idx] = regexp.MustCompile(p)
	}
	return s
}

// SetDelim is a function set charactors as delimiter
func (x *SimpleSplitter) SetDelim(d string) {
	x.delims = d
}

// EnableRegex is disabler of heuristics patterns
func (x *SimpleSplitter) EnableRegex() {
	x.useRegex = true
}

// DisableRegex is disabler of heuristics patterns
func (x *SimpleSplitter) DisableRegex() {
	x.useRegex = false
}

func (x *SimpleSplitter) splitByRegex(chunk *Chunk) []*Chunk {
	if x.useRegex {
		for _, regex := range x.regexList {
			result := regex.FindAllStringIndex(chunk.Data, -1)
			if len(result) > 0 {
				pos := 0
				chunks := make([]*Chunk, len(result)*2+1)
				for idx, m := range result {
					chunks[idx*2] = newChunk(chunk.Data[pos:m[0]])
					chunks[idx*2+1] = newChunk(chunk.Data[m[0]:m[1]])
					chunks[idx*2+1].Freeze = true
					pos = m[1]
				}
				chunks[len(chunks)-1] = newChunk(chunk.Data[pos:])
				return chunks
			}
		}
	}

	res := []*Chunk{chunk}
	return res
}

func (x *SimpleSplitter) splitByDelimiter(chunk *Chunk) []*Chunk {
	var res []*Chunk
	msg := chunk.Data
	for {
		idx := strings.IndexAny(msg, x.delims)
		if idx < 0 {
			if len(msg) > 0 {
				res = append(res, newChunk(msg))
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
			res = append(res, newChunk(s1))
		}

		if len(s2) > 0 {
			// log.Print("add s2: ", s2)
			res = append(res, newChunk(s2))
		}

		msg = s3
		// log.Print("remain: ", msg)
	}

	return res
}

// Split is a function to split log message.
func (x *SimpleSplitter) Split(msg string) []*Chunk {
	chunk := newChunk(msg)
	prevLen := 0
	chunks := []*Chunk{chunk}

	for prevLen != len(chunks) {
		var tmp []*Chunk
		for _, c := range chunks {
			// log.Println(c)
			if c.Freeze {
				tmp = append(tmp, c)
			} else {
				tmp = append(tmp, x.splitByRegex(c)...)
			}
		}
		prevLen = len(chunks)
		chunks = tmp
	}

	var res []*Chunk
	for _, c := range chunks {
		if c.Freeze {
			res = append(res, c)
		} else {
			res = append(res, x.splitByDelimiter(c)...)
		}
	}
	return res
}
