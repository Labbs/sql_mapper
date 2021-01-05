package dump

import (
	"bytes"
	"fmt"
	"log"
	"sql_mapper/config"
	"strings"
	"text/template"

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

var tableTmpl = "--\n" +
	"-- Table structure for table `{{.TableName}}`\n" +
	"--\n" +
	"{{ if .DropTable }}DROP TABLE IF EXISTS `{{.TableName}}`;{{ end }}\n" +
	"CREATE TABLE `{{.TableName}}` (\n" +
	"{{ range $i, $v := .Columns }}" +
	"  `{{$v.Name}}` {{$v.Type}} " +
	"{{ if $v.CollationName }}CHARACTER SET {{ if $v.Charset }}{{$v.Charset}}{{else}}{{$.Charset}}{{end}} COLLATE {{$v.CollationName}} {{ end }}" +
	"{{ if (eq $v.Null \"NO\") }}NOT NULL{{ end }}" +
	"{{ if and (eq $v.Null \"YES\") (ne $v.PrimaryKey true) }}DEFAULT NULL{{ end }} " +
	"{{ if (ne $v.Default \"\") }}DEFAULT {{ if (eq $v.Default \"CURRENT_TIMESTAMP\") }}{{$v.Default}}{{ else }}'{{$v.Default}}'{{ end }}{{ end }}" +
	"{{ if (eq $v.Empty \"1\") }}DEFAULT '' {{ end }}" +
	"{{ if $v.AutoIncrement }}AUTO_INCREMENT{{ end }},\n" +
	"{{ end }}" +
	"{{- if .Primary }}" +
	"PRIMARY KEY ({{range .Primary}}`{{.}}`,{{end}}),\n" +
	"{{ end -}}" +
	"{{- range $i, $v := .Keys }}" +
	"KEY `{{$v.Name}}` ({{range $v.ColumnName}}`{{.}}`,{{end}}),\n" +
	"{{ end -}}" +
	"{{- range $i, $v := .Unique }}" +
	"UNIQUE KEY `{{$v.Name}}` ({{ range $v.Column }}`{{.}}`,{{ end }}),\n" +
	"{{ end -}}" +
	"{{- range $i, $v := .Constraints }}" +
	"{{- if $v.Table }}" +
	"CONSTRAINT `{{$v.Name}}` FOREIGN KEY ({{range $v.Column}}`{{.}}`,{{end}}) REFERENCES `{{$v.Table}}` ({{range $v.TableID}}`{{.}}`,{{end}}){{ if $v.OnDelete}} ON DELETE {{$v.OnDelete}}{{end}}{{ if $v.OnUpdate}} ON UPDATE {{$v.OnUpdate}}{{end}},\n" +
	"{{ end -}}" +
	"{{ end -}}" +
	") ENGINE={{.Engine}} DEFAULT CHARSET={{.Charset}};\n\n"

// ExportTablesSchema ---
func ExportTablesSchema(engine *xorm.Engine) {
	tables, err := engine.Query("select * from information_schema.tables where TABLE_SCHEMA = \"" + config.Dump.Database + "\";")
	if err != nil {
		log.Fatal(err.Error())
	}

	schemata, err := engine.Query("select * from information_schema.SCHEMATA where SCHEMA_NAME = \"" + config.Dump.Database + "\";")
	if err != nil {
		log.Fatal(err.Error())
	}

	charset := string(schemata[0]["DEFAULT_CHARACTER_SET_NAME"])
	collationName := string(schemata[0]["DEFAULT_COLLATION_NAME"])

	fmt.Println("/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;")
	fmt.Println("SET NAMES " + charset + ";\n\n")

	for _, table := range tables {
		var t TableStruct
		columns, err := engine.Query("select *, IF(COLUMN_DEFAULT = '' AND COLUMN_DEFAULT IS NOT NULL, '1', '0') as EMPTY from information_schema.columns where TABLE_SCHEMA = \"" + config.Dump.Database + "\" AND TABLE_NAME = \"" + string(table["TABLE_NAME"]) + "\";")
		if err != nil {
			log.Fatal(err.Error())
		}

		keyColumnUsage, err := engine.Query("select * from information_schema.KEY_COLUMN_USAGE where CONSTRAINT_SCHEMA = \"" + config.Dump.Database + "\" AND CONSTRAINT_NAME <> \"PRIMARY\" AND TABLE_NAME = \"" + string(table["TABLE_NAME"]) + "\";")
		if err != nil {
			log.Fatal(err.Error())
		}

		statistics, err := engine.Query("select * from information_schema.STATISTICS where TABLE_SCHEMA = \"" + config.Dump.Database + "\" AND TABLE_NAME = \"" + string(table["TABLE_NAME"]) + "\" AND NON_UNIQUE = 1;")
		if err != nil {
			log.Fatal(err.Error())
		}

		referentialConstraints, err := engine.Query("select * from information_schema.REFERENTIAL_CONSTRAINTS where UNIQUE_CONSTRAINT_SCHEMA = \"" + config.Dump.Database + "\" AND TABLE_NAME = \"" + string(table["TABLE_NAME"]) + "\";")
		if err != nil {
			log.Fatal(err.Error())
		}

		tableConstraint, err := engine.Query("select * from information_schema.TABLE_CONSTRAINTS where CONSTRAINT_SCHEMA = \"" + config.Dump.Database + "\" AND TABLE_NAME = \"" + string(table["TABLE_NAME"]) + "\";")
		if err != nil {
			log.Fatal(err.Error())
		}

		t.TableName = string(table["TABLE_NAME"])
		t.Engine = string(table["ENGINE"])
		t.Charset = charset
		t.DropTable = true

		for _, column := range columns {
			if string(column["COLUMN_KEY"]) == "PRI" {
				var exist = false
				for _, c := range tableConstraint {
					if string(c["CONSTRAINT_NAME"]) == "PRIMARY" {
						exist = true
					}
				}
				if exist {
					t.Primary = append(t.Primary, string(column["COLUMN_NAME"]))
				}
			}
		}

		for _, column := range columns {
			var c ColumnStruct
			c.Name = string(column["COLUMN_NAME"])
			c.Type = string(column["COLUMN_TYPE"])
			c.Null = string(column["IS_NULLABLE"])
			c.Default = string(column["COLUMN_DEFAULT"])
			c.Empty = string(column["EMPTY"])
			c.ColumnKey = string(column["COLUMN_KEY"])
			if string(column["COLUMN_KEY"]) == "PRI" {
				c.PrimaryKey = true
			}
			if string(column["EXTRA"]) == "auto_increment" {
				c.AutoIncrement = true
			}
			if string(column["CHARACTER_SET_NAME"]) != charset && string(column["CHARACTER_SET_NAME"]) != "" {
				c.Charset = string(column["CHARACTER_SET_NAME"])
			}
			if string(column["COLLATION_NAME"]) != collationName && string(column["CHARACTER_SET_NAME"]) != "" && string(column["COLLATION_NAME"]) != fmt.Sprintf("%s_general_ci", string(column["CHARACTER_SET_NAME"])) {
				c.CollationName = string(column["COLLATION_NAME"])
			}
			t.Columns = append(t.Columns, c)
		}

		for _, col := range keyColumnUsage {
			if string(col["REFERENCED_TABLE_NAME"]) != "" {
				var c ConstraintStruct
				var index int
				var exist = false
				for k, v := range t.Constraints {
					if v.Name == string(col["CONSTRAINT_NAME"]) {
						exist = true
						index = k
					}
				}

				if exist {
					t.Constraints[index].Column = append(t.Constraints[index].Column, string(col["COLUMN_NAME"]))
					t.Constraints[index].TableID = append(t.Constraints[index].TableID, string(col["REFERENCED_COLUMN_NAME"]))
				} else {
					c.Name = string(col["CONSTRAINT_NAME"])
					c.Column = append(c.Column, string(col["COLUMN_NAME"]))
					c.Table = string(col["REFERENCED_TABLE_NAME"])
					c.TableID = append(c.TableID, string(col["REFERENCED_COLUMN_NAME"]))
					for _, d := range referentialConstraints {
						if string(col["CONSTRAINT_NAME"]) == string(d["CONSTRAINT_NAME"]) && string(d["UPDATE_RULE"]) != "RESTRICT" {
							c.OnUpdate = string(d["UPDATE_RULE"])
						}
						if string(col["CONSTRAINT_NAME"]) == string(d["CONSTRAINT_NAME"]) && string(d["DELETE_RULE"]) != "RESTRICT" {
							c.OnDelete = string(d["DELETE_RULE"])
						}
					}
					t.Constraints = append(t.Constraints, c)
				}

			}
		}

		for _, col := range statistics {
			var c KeyStruct
			var index int
			var exist = false
			for k, v := range t.Keys {
				if v.Name == string(col["INDEX_NAME"]) {
					exist = true
					index = k
				}
			}
			if exist {
				t.Keys[index].ColumnName = append(t.Keys[index].ColumnName, string(col["COLUMN_NAME"]))
			} else {
				c.Name = string(col["INDEX_NAME"])
				c.ColumnName = append(c.ColumnName, string(col["COLUMN_NAME"]))
				t.Keys = append(t.Keys, c)
			}
		}

		for _, col := range keyColumnUsage {
			if string(col["REFERENCED_TABLE_NAME"]) == "" {
				var index int
				var exist = false
				for k, v := range t.Unique {
					if v.Name == string(col["CONSTRAINT_NAME"]) {
						exist = true
						index = k
					}
				}
				if exist {
					t.Unique[index].Column = append(t.Unique[index].Column, string(col["COLUMN_NAME"]))
				} else {
					var c UniqueStruct
					c.Name = string(col["CONSTRAINT_NAME"])
					c.Column = append(c.Column, string(col["COLUMN_NAME"]))
					t.Unique = append(t.Unique, c)
				}
			}
		}

		var tpl bytes.Buffer

		_t := template.Must(template.New("table").Parse(tableTmpl))
		err = _t.Execute(&tpl, t)
		if err != nil {
			panic(err)
		}
		t1 := strings.Replace(tpl.String(), ",\n)", "\n)", -1)
		fmt.Printf("%s", strings.Replace(t1, ",)", ")", -1))
	}
}
