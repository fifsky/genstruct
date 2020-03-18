package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql/v2"

	"github.com/fifsky/genstruct/generator"
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

	gen := generator.NewGenerator(gosql.Use("default"))

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
				cmd, err := generator.GetParams(cmds, 1)
				if err != nil {
					return err
				}
				err = link(cmd)
				if err == nil {
					fmt.Println("Database changed")
				}
				return err

			case "g":
				cmd, err := generator.GetParams(cmds, 1)
				if err != nil {
					return err
				}

				tag, _ := generator.GetParams(cmds, 2)
				tags := strings.Split(tag, ",")
				if len(tags) == 0 {
					tags = []string{"db", "json"}
				}

				out, err := gen.ShowStruct(cmd, tags)
				if err != nil {
					return err
				}
				fmt.Println(string(out))
			case "exit":
				fmt.Println("Bye!")
				os.Exit(0)
			default:
				start := time.Now()
				datas, err := gen.Exec(line)
				if err != nil {
					return err
				}
				gen.ShowTable(datas, start)
			}
			return
		}()
	}
}
