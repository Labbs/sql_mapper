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
	Diff struct {
		Engine string

		SourceHostname string
		SourceUsername string
		SourcePassword string
		SourceDatabase string
		SourcePort     int

		DestHostname string
		DestUsername string
		DestPassword string
		DestDatabase string
		DestPort     int
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
	}
}

func DiffSchema() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "source_hostname",
			Value:       "localhost",
			Destination: &Diff.SourceHostname,
			EnvVar:      "SOURCE_HOSTNAME",
		},
		cli.StringFlag{
			Name:        "source_username",
			Value:       "root",
			Destination: &Diff.SourceUsername,
			EnvVar:      "SOURCE_USERNAME",
		},
		cli.StringFlag{
			Name:        "source_password",
			Value:       "root",
			Destination: &Diff.SourcePassword,
			EnvVar:      "SOURCE_PASSWORD",
		},
		cli.StringFlag{
			Name:        "source_database",
			Value:       "mysql",
			Destination: &Diff.SourceDatabase,
			EnvVar:      "SOURCE_DATABASE",
		},
		cli.IntFlag{
			Name:        "source_port",
			Value:       3306,
			Destination: &Diff.SourcePort,
			EnvVar:      "SOURCE_PORT",
		},
		cli.StringFlag{
			Name:        "dest_hostname",
			Value:       "localhost",
			Destination: &Diff.DestHostname,
			EnvVar:      "DEST_HOSTNAME",
		},
		cli.StringFlag{
			Name:        "dest_username",
			Value:       "root",
			Destination: &Diff.DestUsername,
			EnvVar:      "DEST_USERNAME",
		},
		cli.StringFlag{
			Name:        "dest_password",
			Value:       "root",
			Destination: &Diff.DestPassword,
			EnvVar:      "DEST_PASSWORD",
		},
		cli.StringFlag{
			Name:        "dest_database",
			Value:       "mysql",
			Destination: &Diff.DestDatabase,
			EnvVar:      "DEST_DATABASE",
		},
		cli.IntFlag{
			Name:        "dest_port",
			Value:       3306,
			Destination: &Diff.DestPort,
			EnvVar:      "DEST_PORT",
		},
		cli.StringFlag{
			Name:        "engine",
			Value:       "mysql",
			Destination: &Diff.Engine,
			EnvVar:      "ENGINE",
		},
	}
}
