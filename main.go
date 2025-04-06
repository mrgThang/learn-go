package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	migrateV4 "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/mrgThang/learn-go/config"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(_ []string) error {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "learn-go"
	app.Commands = []*cli.Command{
		{
			Name:        "migrate",
			Usage:       "doing migrations",
			Subcommands: doMigration(cfg.MigrationFolder, cfg.Mysql.String()),
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}

	return nil
}

func doMigration(sourceUrl string, databaseUrl string) []*cli.Command {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return []*cli.Command{
		{
			Name:  "up",
			Usage: "migrate up",
			Action: func(c *cli.Context) error {
				m, err := migrateV4.New(sourceUrl, databaseUrl)
				if err != nil {
					logger.Fatal("Error create migrations", zap.Error(err))
				}

				logger.Info("migrate up")
				if err = m.Up(); err != nil && !errors.Is(err, migrateV4.ErrNoChange) {
					logger.Error(err.Error())
					return err
				}
				return nil
			},
		},
		{
			Name:  "down",
			Usage: "migrate down",
			Action: func(c *cli.Context) error {
				m, err := migrateV4.New(sourceUrl, databaseUrl)
				if err != nil {
					logger.Fatal("Error create migrations", zap.Error(err))
				}

				down, err := strconv.Atoi(c.Args().First())
				if err != nil {
					logger.Fatal("rev should be a number", zap.Error(err))
				}

				logger.Info("migrate down", zap.Int("down", -down))
				if err = m.Steps(-down); err != nil {
					logger.Error(err.Error())
					return err
				}

				return nil
			},
		},
	}
}
