package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/urfave/cli"

	"sql_mapper/config"
	"sql_mapper/diff"
	"sql_mapper/dump"
	"sql_mapper/render"
)

var (
	engine       *xorm.Engine
	sourceEngine *xorm.Engine
	destEngine   *xorm.Engine
)

var version = "v1.0.0"

func main() {
	app := cli.NewApp()

	app.Name = "sql_dumper"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:   "schema",
			Usage:  "dump mysql schema",
			Flags:  config.DumpSchema(),
			Action: dumpSchema,
		},
		{
			Name:   "render",
			Usage:  "render an image from database schema",
			Flags:  config.DumpSchema(),
			Action: renderImage,
		},
		{
			Name:   "diff",
			Usage:  "difference schema between two databases",
			Flags:  config.DiffSchema(),
			Action: diffSchema,
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

	dump.ExportTablesSchema(engine)
}

func renderImage(ctx *cli.Context) {
	var err error
	engine, err = xorm.NewEngine(config.Dump.Engine, fmt.Sprintf("%s:%s@(%s:%d)/information_schema", config.Dump.Username, config.Dump.Password, config.Dump.Hostname, config.Dump.Port))
	engine.ShowSQL(false)

	if err != nil {
		log.Fatal(err.Error())
	}

	render.RenderSchema(engine)
}

func diffSchema(ctx *cli.Context) {
	var err error

	sourceEngine, err = xorm.NewEngine(config.Diff.Engine, fmt.Sprintf("%s:%s@(%s:%d)/information_schema", config.Diff.SourceUsername, config.Diff.SourcePassword, config.Diff.SourceHostname, config.Diff.SourcePort))
	sourceEngine.ShowSQL(false)

	if err != nil {
		log.Fatal(err.Error())
	}

	destEngine, err = xorm.NewEngine(config.Diff.Engine, fmt.Sprintf("%s:%s@(%s:%d)/information_schema", config.Diff.DestUsername, config.Diff.DestPassword, config.Diff.DestHostname, config.Diff.DestPort))
	destEngine.ShowSQL(false)

	if err != nil {
		log.Fatal(err.Error())
	}

	diff.DiffSchema(sourceEngine, destEngine)
}
