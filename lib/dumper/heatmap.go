package dumper

import (
	"fmt"
	logptn "github.com/m-mizutani/logptn/lib"
	"html/template"
	// logger "log"
)

type HeatmapDumper struct {
	fileDumper
}

// NewHeatmapDumper is Constructor of Heatmapdumper
func NewHeatmapDumper(fpath string) (*HeatmapDumper, error) {
	dumper := HeatmapDumper{}
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
td { padding: 2px }
thead td.sample { text-align: center; }
thead td.cell {  -webkit-writing-mode: vertical-rl; -ms-writing-mode: tb-rl; writing-mode: vertical-rl; }
td.sample { width: 840px; border: 1px solid #999; }
td.cell { width: 24px; font-size: 12px; text-align: center; border: 1px solid #aaa; }
td.total { text-align: right; border: 1px solid #999; padding: 4px; }
span.var { color: #D04255 }
  --></style>  
</head>
<body>
<table>
<thead>
<tr>
<td class="sample">Log formats</td>
{{range .IndexList}}<td class="cell">{{.}}</td>{{end}}
<td class="total">Total</td>
</tr>
</thead>
<tbody>
{{range .Formats}}
  <tr>
  <td class="sample">{{range .Segments}}
    {{if .IsVar}}<span class="var">*</span>{{else}}<span>{{.Word | html}}{{end}}</span>
  {{end}}</td>
  {{range .Cells}}<td class="cell">{{if .Has }}{{.Count}}{{end}}</td>{{end}}
  <td class="total">{{.Total}}</td>
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

type mtxRenderData struct {
	Formats   []*mtxFormat
	IndexList []string
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
func (x *HeatmapDumper) DumpFormat(formats []*logptn.Format) error {
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

	renderData := mtxRenderData{}
	for i := 0; i < total; i = i + unitSize {
		renderData.IndexList = append(renderData.IndexList,
			fmt.Sprintf("%d - %d", i+1, i+unitSize))
	}

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

	renderData.Formats = mf
	err = t.Execute(x.out, renderData)
	if err != nil {
		return err
	}
	return nil
}

func (x *HeatmapDumper) DumpLog(logs []*logptn.Log) error {
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}
