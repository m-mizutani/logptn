package logptn

type Chunk struct {
	chunk string
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
