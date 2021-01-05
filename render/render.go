package render

import (
	"fmt"
	"log"
	"sql_mapper/config"

	"github.com/awalterschulze/gographviz"
	"github.com/go-xorm/xorm"
)

type TableStruct struct {
	DropTable   bool
	TableName   string
	Engine      string
	Charset     string
	Primary     []string
	Columns     []ColumnStruct
	Constraints []ConstraintStruct
	Unique      []UniqueStruct
	Keys        []KeyStruct
}

type ColumnStruct struct {
	Name          string
	Type          string
	Null          string
	Default       string
	ColumnKey     string
	PrimaryKey    bool
	AutoIncrement bool
	Charset       string
	CollationName string
}

type KeyStruct struct {
	Name       string
	ColumnName []string
}

type ConstraintStruct struct {
	Name     string
	Column   []string
	Table    string
	TableID  []string
	OnDelete string
	OnUpdate string
}

type UniqueStruct struct {
	Name   string
	Column []string
}

type ColumnTypeStruct struct {
	Name      string
	TableName string
	Type      string
}

var tmpl = "digraph G {" +

	"}"

func RenderSchema(engine *xorm.Engine) {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}

	tables, err := engine.Query("select * from information_schema.tables where TABLE_SCHEMA = \"" + config.Dump.Database + "\";")
	if err != nil {
		log.Fatal(err.Error())
	}

	// schemata, err := engine.Query("select * from information_schema.SCHEMATA where SCHEMA_NAME = \"" + config.Dump.Database + "\";")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// charset := string(schemata[0]["DEFAULT_CHARACTER_SET_NAME"])
	// collationName := string(schemata[0]["DEFAULT_COLLATION_NAME"])

	for _, table := range tables {
		graph.AddNode("G", string(table["TABLE_NAME"]), nil)
	}

	for _, table := range tables {
		keyColumnUsage, err := engine.Query("select * from information_schema.KEY_COLUMN_USAGE where CONSTRAINT_SCHEMA = \"" + config.Dump.Database + "\" AND CONSTRAINT_NAME <> \"PRIMARY\" AND TABLE_NAME = \"" + string(table["TABLE_NAME"]) + "\";")
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, col := range keyColumnUsage {
			if string(col["REFERENCED_TABLE_NAME"]) != "" {
				graph.AddPortEdge(string(table["TABLE_NAME"]), string(col["COLUMN_NAME"]), string(col["REFERENCED_TABLE_NAME"]), string(col["REFERENCED_COLUMN_NAME"]), true, nil)
			}
		}
	}

	output := graph.String()
	fmt.Println(output)
}
