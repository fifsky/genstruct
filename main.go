package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql"

	"github.com/fifsky/genstruct/lib"
)

var (
	host     = flag.String("h", "localhost", "database host")
	user     = flag.String("u", "root", "database user")
	password = flag.String("P", "", "database passwrd")
	port     = flag.String("p", "3306", "database port")
)

func link(database string) error {
	configs := make(map[string]*gosql.Config)
	configs["default"] = &gosql.Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *user, *password, *host, *port, database) + "?charset=utf8&parseTime=True&loc=Asia%2FShanghai",
		ShowSql: false,
	}
	return gosql.Connect(configs)
}

func main() {
	flag.Parse()
	gosql.FatalExit = false
	err := link("")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for input.Scan() {
		func() (err error) {
			defer func() {
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print("> ")
			}()

			line := strings.TrimRight(strings.TrimSpace(input.Text()), ";")

			if line == "" {
				return
			}

			cmds := strings.Split(line, " ")

			switch cmds[0] {
			case "use":
				cmd, err := lib.GetParams(cmds, 1)
				if err != nil {
					return err
				}
				err = link(cmd)
				if err == nil {
					fmt.Println("Database changed")
				}
				return err

			case "g":
				cmd, err := lib.GetParams(cmds, 1)
				if err != nil {
					return err
				}

				tag, _ := lib.GetParams(cmds, 2)

				err = lib.ShowStruct(cmd, tag)
				if err != nil {
					return err
				}
			case "exit":
				fmt.Println("Bye!")
				os.Exit(0)
			default:
				start := time.Now()
				datas, err := lib.Exec(line)
				if err != nil {
					return err
				}
				lib.ShowTable(datas, start)
			}
			return
		}()
	}
}
