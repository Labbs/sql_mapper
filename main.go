package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/urfave/cli"

	"sql_dumper/config"
	"sql_dumper/dump"
)

var engine *xorm.Engine

var version = "v1.0.0"

func main() {
	app := cli.NewApp()

	app.Name = "sql_dumper"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:   "dump",
			Usage:  "dump mysql",
			Flags:  config.DumpSchema(),
			Action: dumpSchema,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func dumpSchema(ctx *cli.Context) {
	var err error
	engine, err = xorm.NewEngine(config.Dump.Engine, fmt.Sprintf("%s:%s@(%s:%d)/information_schema", config.Dump.Username, config.Dump.Password, config.Dump.Hostname, config.Dump.Port))
	engine.ShowSQL(false)

	if err != nil {
		log.Fatal(err.Error())
	}

	dump.CreateTable(engine)
}
