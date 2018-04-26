package logptn

type Chunk struct {
	Data string
}

type Log struct {
	text  string
	chunk []Chunk
}

func NewLog(line string) *Log {
	log := Log{}
	log.text = line
	return &log
}

func NewChunk(d string) *Chunk {
	c := Chunk{}
	c.Data = d
	return &c
}
