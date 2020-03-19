package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"
	"time"

	"github.com/ilibs/gosql/v2"
	"github.com/olekukonko/tablewriter"
)

const tmplContent = `
package {{ .TableName }}
{{ if .ExistTime }}
import (
	"time"
)
{{end}}
type {{ .StructName }} struct {
    {{ range $i,$v := .Columns }}{{ .StructField }}    {{ .Type }}    ` + "\u0060" + `{{ range $j,$tag := $.OtherTags }} {{ $tag }}:"{{ $v.Field }}"{{ end }}` + "\u0060" + `{{ if ne .Comment "" }} // {{.Comment}}{{ end }}{{ if ne $i $.Len }}` + "\n" + `{{ end }}{{ end }}
}

func ({{ .ShortName }} *{{ .StructName }}) TableName() string {
    return "{{ .TableName }}"
}

func ({{ .ShortName }} *{{ .StructName }}) PK() string {
    return "{{ .PrimaryKey }}"
}
`

type Attr struct {
	StructField string
	Field       string
	Type        string
	Comment     string
}

type TableInfo struct {
	Columns    []*Attr
	Len        int
	OtherTags  []string
	TableName  string
	ShortName  string
	StructName string
	Database   string
	PrimaryKey string
	ExistTime  bool
}

type Generator struct {
	db *gosql.DB
}

func NewGenerator(db *gosql.DB) *Generator {
	return &Generator{db: db}
}

func (g *Generator) Exec(query string) ([]map[string]interface{}, error) {
	rows, err := gosql.Queryx(query)
	if err != nil {
		return nil, err
	}

	var datas []map[string]interface{}

	for rows.Next() {
		data := make(map[string]interface{})
		rows.MapScan(data)
		datas = append(datas, data)
	}
	return datas, nil
}

func (g *Generator) ShowTable(datas []map[string]interface{}, start time.Time) {
	if len(datas) > 0 {
		header, cells := formatTable(datas)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		table.AppendBulk(cells)
		table.Render()
		end := time.Now()
		fmt.Println(fmt.Sprintf("%d rows in set (%.2f sec)", len(cells), float64(end.UnixNano()-start.UnixNano())/1e9))
	} else {
		fmt.Println("No Result")
	}
}

func (g *Generator) ShowStruct(table string, tags []string) ([]byte, error) {
	query := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", table)
	datas, err := g.Exec(query)
	if err != nil {
		return nil, err
	}

	var databaseName string
	err = gosql.QueryRowx("select database()").Scan(&databaseName)

	if err != nil {
		return nil, err
	}

	info := &TableInfo{
		OtherTags:  tags,
		Columns:    make([]*Attr, 0),
		TableName:  table,
		ShortName:  table[0:1],
		StructName: titleCasedName(table),
		Database:   databaseName,
	}

	var existTime = false
	for _, v := range datas {
		m := mapToString(v)
		tp := typeFormat(m["Type"], m["Null"])

		if tp == "time.Time" {
			existTime = true
		}

		attr := &Attr{
			StructField: titleCasedName(m["Field"]),
			Field:       m["Field"],
			Type:        tp,
			Comment:     m["Comment"],
		}

		info.Columns = append(info.Columns, attr)
		if m["Key"] == "PRI" {
			info.PrimaryKey = attr.Field
		}
	}

	info.ExistTime = existTime

	info.Len = len(info.Columns) - 1

	tmpl, err := template.New("struct").Parse(tmplContent)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, info)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
