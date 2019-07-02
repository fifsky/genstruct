package lib

import (
	"fmt"
	"os"
	"time"

	"github.com/ilibs/gosql"
	"github.com/olekukonko/tablewriter"
)

func Exec(query string) ([]map[string]interface{}, error) {
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

func ShowTable(datas []map[string]interface{}, start time.Time) {
	if len(datas) > 0 {
		header, cells := formatTable(datas)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		table.AppendBulk(cells)
		table.Render()
		end := time.Now()
		fmt.Println(fmt.Sprintf("%d rows in set (%.2f sec)", len(cells), float64(end.UnixNano()-start.UnixNano())/float64(1e9)))
	} else {
		fmt.Println("No Result")
	}
}
