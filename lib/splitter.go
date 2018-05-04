package logptn

import (
	// "log"
	"regexp"
	"strings"
)

type Splitter struct {
	delims    string
	regexList []*regexp.Regexp
}

func NewSplitter() *Splitter {
	s := Splitter{}
	s.delims = " \t!,:;[]{}()<>=|\\*\"'"
	heuristicsPatterns := []string{
		`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+`,                           // DateTime
		`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`,                              // DateTime
		`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`,                               // DateTime
		`\d{4}/\d{2}/\d{2}`,                                                 // Date
		`\d{4}-\d{2}-\d{2}`,                                                 // Date
		`\d{2}:\d{2}:\d{2}.\d+`,                                             // Time
		`\d{2}:\d{2}:\d{2}`,                                                 // Time
		`[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*`, // Mail address
		`(\d{1,3}\.){3}\d{1,3}`,                                             // IPv4 address
	}

	s.regexList = make([]*regexp.Regexp, len(heuristicsPatterns))
	for idx, p := range heuristicsPatterns {
		s.regexList[idx] = regexp.MustCompile(p)
	}
	return &s
}

func (x *Splitter) SetDelim(d string) {
	x.delims = d
}

func (x *Splitter) splitByRegex(chunk *Chunk) []*Chunk {
	for _, regex := range x.regexList {
		result := regex.FindAllStringIndex(chunk.Data, -1)
		if len(result) > 0 {
			pos := 0
			chunks := make([]*Chunk, len(result)*2+1)
			for idx, m := range result {
				chunks[idx*2] = NewChunk(chunk.Data[pos:m[0]])
				chunks[idx*2+1] = NewChunk(chunk.Data[m[0]:m[1]])
				chunks[idx*2+1].Freeze = true
				pos = m[1]
			}
			chunks[len(chunks)-1] = NewChunk(chunk.Data[pos:])
			return chunks
		}
	}
	res := []*Chunk{chunk}
	return res
}

func (x *Splitter) splitByDelimiter(chunk *Chunk) []*Chunk {
	var res []*Chunk
	msg := chunk.Data
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

func (x *Splitter) Split(msg string) []*Chunk {
	chunk := NewChunk(msg)
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
