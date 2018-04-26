package logptn

type Chunk struct {
	Data string
}

type Log struct {
	text  string
	Chunk []*Chunk
}

func NewLog(line string, sp *Splitter) *Log {
	log := Log{}
	log.text = line
	log.Chunk = sp.Split(line)
	return &log
}

func NewChunk(d string) *Chunk {
	c := Chunk{}
	c.Data = d
	return &c
}
