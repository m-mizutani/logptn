package logptn

type Chunk struct {
	Data   string
	Freeze bool
}

type Log struct {
	text   string
	Chunk  []*Chunk
	format *Format
}

func newLog(line string, sp Splitter) *Log {
	log := Log{}
	log.text = line
	log.Chunk = sp.Split(line)
	return &log
}

func newChunk(d string) *Chunk {
	c := Chunk{}
	c.Data = d
	c.Freeze = false
	return &c
}

func (x *Chunk) Clone() *Chunk {
	c := newChunk(x.Data)
	c.Freeze = x.Freeze
	return c
}

func (x *Chunk) String() string {
	return x.Data
}
