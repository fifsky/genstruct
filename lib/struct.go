package lib

import (
	"fmt"
	"text/template"
	"os"
	"github.com/ilibs/gosql"
	"os/exec"
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
	ExistTime  bool
}

var tmplContent = `
package {{ .TableName }}
{{ if .ExistTime }}
import (
	"time"
)
{{end}}
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

	var existTime = false
	for _, v := range datas {
		m := mapToString(v)
		tp := typeFormat(m["Type"])

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
		return err
	}

	tmpFile := "/tmp/genstruct"

	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	err = tmpl.Execute(f, info)

	c := exec.Command("gofmt", "-s", tmpFile)
	c.Dir = "/tmp"
	out, err := c.CombinedOutput()

	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
