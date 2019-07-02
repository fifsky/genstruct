package lib

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql"
)

func TestMain(m *testing.M) {
	configs := make(map[string]*gosql.Config)

	dsn := os.Getenv("MYSQL_TEST_DSN")

	if dsn == "" {
		dsn = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	}

	configs["default"] = &gosql.Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     dsn,
		ShowSql: true,
	}

	gosql.Connect(configs)

	m.Run()
}
