package genlogfmt

type Log struct {
	text string
}

type Generator struct {
	logs []Log
}

func (x *Generator) ReadFile(fpath string) error {
	return nil
}
