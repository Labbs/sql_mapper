package diff

import (
	"log"
	"sql_mapper/config"

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
	Empty         string
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

func DiffSchema(sourceEngine *xorm.Engine, destEngine *xorm.Engine) {
	sourceTables, err := sourceEngine.Query("select * from information_schema.tables where TABLE_SCHEMA = \"" + config.Diff.SourceDatabase + "\";")
	if err != nil {
		log.Fatal(err.Error())
	}

	// destTables, err := destEngine.Query("select * from information_schema.tables where TABLE_SCHEMA = \"" + config.Diff.DestDatabase + "\";")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// sourceSchemata, err := sourceEngine.Query("select * from information_schema.SCHEMATA where SCHEMA_NAME = \"" + config.Diff.SourceDatabase + "\";")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// destSchemata, err := destEngine.Query("select * from information_schema.SCHEMATA where SCHEMA_NAME = \"" + config.Diff.DestDatabase + "\";")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// var notExist []string
	// for _, sTable := range sourceTables {
	// 	exist := false
	// 	for _, dTable := range destTables {
	// 		if string(sTable["TABLE_NAME"]) == string(dTable["TABLE_NAME"]) {
	// 			exist = true
	// 		}
	// 	}
	// 	if ! exist {
	// 		notExist = append(notExist, string(sTable["TABLE_NAME"]))
	// 		// print table schema
	// 	}
	// }

	// for _, sTable := range sourceTables {
	// 	if ! helpers.ExistInArray(string(sTable["TABLE_NAME"]), notExist) {

	// 	}
	// }
}
