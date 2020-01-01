package config

import "github.com/urfave/cli"

var (
	Dump struct {
		Hostname   string
		Username   string
		Password   string
		Database   string
		Port       int
		Engine     string
		SchemaOnly bool
	}
)

func DumpSchema() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "host, hostname",
			Value:       "localhost",
			Destination: &Dump.Hostname,
			EnvVar:      "HOSTNAME",
		},
		cli.StringFlag{
			Name:        "user, username",
			Value:       "root",
			Destination: &Dump.Username,
			EnvVar:      "USERNAME",
		},
		cli.StringFlag{
			Name:        "pass, password",
			Value:       "root",
			Destination: &Dump.Password,
			EnvVar:      "PASSWORD",
		},
		cli.StringFlag{
			Name:        "database",
			Value:       "mysql",
			Destination: &Dump.Database,
			EnvVar:      "DATABASE",
		},
		cli.IntFlag{
			Name:        "port",
			Value:       3306,
			Destination: &Dump.Port,
			EnvVar:      "PORT",
		},
		cli.StringFlag{
			Name:        "engine",
			Value:       "mysql",
			Destination: &Dump.Engine,
			EnvVar:      "ENGINE",
		},
		cli.BoolFlag{
			Name:        "schema_only",
			Destination: &Dump.SchemaOnly,
			EnvVar:      "SCHEMA_ONLY",
		},
	}
}
