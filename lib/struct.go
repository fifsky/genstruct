package lib

import (
	"fmt"
	"text/template"
	"os"
	"github.com/ilibs/gosql"
)

type Attr struct {
	StructField string
	Field       string
	Type        string
	Comment     string
}

type TableInfo struct {
	Columns    []*Attr
	Len        int
	TableName  string
	StructName string
	Database   string
	PrimaryKey string
}

var tmplContent = `
type {{ .StructName }} struct {
    {{ range $i,$v := .Columns }}{{ .StructField }}    {{ .Type }}    ` + "\u0060" + `json:"{{ .Field }}" db:"{{ .Field }}"` + "\u0060{{ if ne $i $.Len }}\n    " + `{{ end }}{{ end }}
}

func (this *{{ .StructName }}) DbName() string {
    return {{ .Database }}
}

func (this *{{ .StructName }}) TableName() string {
    return {{ .TableName }}
}

func (this *{{ .StructName }}) PK() string {
    return {{ .PrimaryKey }}
}
`

func ShowStruct(cmd string) error {
	query := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", cmd)
	datas, err := Exec(query)
	if err != nil {
		return err
	}

	var databaseName string
	err = gosql.QueryRowx("select database()").Scan(&databaseName)

	if err != nil {
		return err
	}

	info := &TableInfo{
		Columns:    make([]*Attr, 0),
		TableName:  cmd,
		StructName: titleCasedName(cmd),
		Database:   databaseName,
	}

	for _, v := range datas {
		m := mapToString(v)

		attr := &Attr{
			StructField: titleCasedName(m["Field"]),
			Field:       m["Field"],
			Type:        typeFormat(m["Type"]),
			Comment:     m["Comment"],
		}

		info.Columns = append(info.Columns, attr)
		if m["Key"] == "PRI" {
			info.PrimaryKey = attr.Field
		}
	}

	info.Len = len(info.Columns) - 1

	tmpl, err := template.New("struct").Parse(tmplContent)
	if err != nil {
		return err
	}
	err = tmpl.Execute(os.Stdout, info)

	if err != nil {
		return err
	}

	return nil
}
