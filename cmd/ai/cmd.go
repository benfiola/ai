package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"

	"github.com/benfiola/ai/pkg/core"
	"github.com/benfiola/ai/pkg/db"
	"github.com/benfiola/ai/pkg/server"
	"github.com/benfiola/ai/pkg/version"
	"github.com/urfave/cli/v3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))

	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintf(cmd.Root().Writer, "%s", cmd.Root().Version)
	}
	err := (&cli.Command{
		Version: version.Get(),
		Commands: []*cli.Command{
			{
				Name:        "migrate",
				Description: "migrates the database",
				Action: func(ctx context.Context, c *cli.Command) error {
					url, err := url.Parse(c.String("database-url"))
					if err != nil {
						return err
					}

					database, err := db.New(db.Opts{
						Logger: logger.With("component", "db"),
						URL:    url,
					})
					if err != nil {
						return err
					}

					versionStr := c.StringArg("version")
					var version db.Version
					switch versionStr {
					case "head":
						version = db.VersionHead
					case "base":
						version = db.VersionBase
					default:
						versionInt, err := strconv.Atoi(versionStr)
						if err != nil {
							break
						}
						version = db.Version(versionInt)
					}
					if err != nil {
						return err
					}

					return database.Migrate(version)
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "version",
						UsageText: "version to migrate database to",
					},
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "database-url",
						Usage:   "database to apply migrations to",
						Sources: cli.EnvVars("DATABASE_URL"),
					},
				},
			},
			{
				Name:        "server",
				Description: "starts the ai server",
				Action: func(ctx context.Context, c *cli.Command) error {
					databaseUrl, err := url.Parse(c.String("database-url"))
					if err != nil {
						return err
					}

					db, err := db.New(db.Opts{
						Logger: logger.With("component", "db"),
						URL:    databaseUrl,
					})
					if err != nil {
						return err
					}

					core, err := core.New(core.Opts{
						DB:     db,
						Logger: logger.With("component", "core"),
					})
					if err != nil {
						return err
					}

					server, err := server.New(server.Opts{
						Core:   core,
						Logger: logger.With("component", "server"),
					})
					if err != nil {
						return err
					}

					return server.Run(ctx)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "database-url",
						Usage:   "database to apply migrations to",
						Sources: cli.EnvVars("DATABASE_URL"),
					},
				},
			},
		},
	}).Run(context.Background(), os.Args)

	code := 0
	if err != nil {
		logger.Error("command failed", "error", err.Error())
		code = 1
	}

	os.Exit(code)
}
