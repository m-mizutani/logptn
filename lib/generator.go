package logptn

type Generator struct {
	logs    []*Log
	formats []*Format
}

func (x *Generator) ReadFile(fpath string) error {
	return nil
}

func (x *Generator) ReadLine(log string) error {
	return nil
}

func (x *Generator) End() {
	clusters := Clustering(x.logs)
	for _, cluster := range clusters {
		format := GenFormat(&cluster)
		x.formats = append(x.formats, format)
	}
}

func (x *Generator) Formats() []*Format {
	return x.formats
}
