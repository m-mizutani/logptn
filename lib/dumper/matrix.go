package dumper

import (
	"fmt"
	logptn "github.com/m-mizutani/logptn/lib"
	"html/template"
	// logger "log"
)

type MatrixDumper struct {
	fileDumper
}

// NewMatrixDumper is Constructor of Matrixdumper
func NewMatrixDumper(fpath string) (*MatrixDumper, error) {
	dumper := MatrixDumper{}
	if err := dumper.open(fpath); err != nil {
		return nil, err
	}
	return &dumper, nil
}

func htmlTemplate() string {
	tpl := `<!DOCTYPE html>
<html>
<head>
  <style type="text/css"><!--
* { font-family: 'Helvetica Neue', sans-serif; font-size: 14px; }
table { border-collapse: collapse; border-spacing: 0; }
td.sample { width: 840px }
td.cell { width: 24px; font-size: 12px; text-align: center; border: 1px solid #999; }
span.var { color: #D04255 }
  --></style>  
</head>
<body>
<table>
<tbody>
{{range .}}
  <tr>
  <td class="sample">{{range .Segments}}
    {{if .IsVar}}<span class="var">*</span>{{else}}<span>{{.Word | html}}{{end}}</span>
  {{end}}</td>
  {{range .Cells}}<td class="cell">{{if .Has }}{{.Count}}{{end}}</td>{{end}}
  <td>{{.Total}}</td>
  </tr>
{{end}}
</tbody>
</table>
</body>
</html>
`
	return tpl
}

type mtxSegment struct {
	Word  string
	IsVar bool
}

type mtxCell struct {
	Count int
	Ratio float64
	Has   bool
}

type mtxFormat struct {
	Segments []*mtxSegment
	Cells    []mtxCell
	Total    int
	original *logptn.Format
}

func calcUnitSize(total, width int) int {
	if total <= width {
		return 1
	}

	baseSize := total / width
	p := 1
	for {
		if baseSize < 10 {
			return p * baseSize
		}
		baseSize = baseSize / 10
		p = p * 10
	}
}

// DumpFormat
func (x *MatrixDumper) DumpFormat(formats []*logptn.Format) error {
	const width = 20
	mf := []*mtxFormat{}
	total := 0

	for _, format := range formats {
		total += len(format.Cluster.Logs())

		var ptr *mtxSegment
		f := mtxFormat{}
		f.original = format
		for _, seg := range format.Segments {
			if seg.Fixed() {
				if ptr == nil {
					ptr = &mtxSegment{Word: seg.Text(), IsVar: false}
					f.Segments = append(f.Segments, ptr)
				} else {
					ptr.Word += seg.Text()
				}
			} else {
				f.Segments = append(f.Segments, &mtxSegment{Word: "*", IsVar: true})
				ptr = nil
			}
		}

		mf = append(mf, &f)
	}

	// logger.Println("total = ", total)
	unitSize := calcUnitSize(total, width)
	// logger.Println("utniSize = ", unitSize)
	actualWidth := (total / unitSize) + 1
	// logger.Println("actWidth =", actualWidth)

	for _, f := range mf {
		f.Cells = make([]mtxCell, actualWidth)
		f.Total = len(f.original.Cluster.Logs())
		for _, log := range f.original.Cluster.Logs() {
			idx := log.Index / unitSize
			f.Cells[idx].Count++
			f.Cells[idx].Has = true
		}
	}

	t, err := template.New("test").Parse(htmlTemplate())
	if err != nil {
		return err
	}

	err = t.Execute(x.out, mf)
	if err != nil {
		return err
	}
	return nil
}

func (x *MatrixDumper) DumpLog(logs []*logptn.Log) error {
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}
