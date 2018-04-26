package logptn

type Chunk struct {
	chunk string
}

type Log struct {
	text  string
	chunk []Chunk
}
