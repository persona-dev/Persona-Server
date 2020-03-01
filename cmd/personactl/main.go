package main

import (
	"database/sql"
	"fmt"
	"os"

	config "github.com/eniehack/persona-server/configs"
	migrate "github.com/rubenv/sql-migrate"
	cli "github.com/urfave/cli/v2"
)

type Gateway struct {
	db        *sql.DB
	migration *migrate.FileMigrationSource
}

type Repository interface {
	ApplyMigrationsToLatest() (int, error)
	RollbackMigrations(applyVersions uint) (int, error)
	ApplyMigrations(applyVersions uint) (int, error)
}

func NewGateway(db *sql.DB, migration *migrate.FileMigrationSource) Repository {
	return &Gateway{
		db:        db,
		migration: migration,
	}
}

func (gateway *Gateway) ApplyMigrationsToLatest() (int, error) {
	appliedmigrationversion, err := migrate.Exec(gateway.db, "postgres", *gateway.migration, migrate.Up)
	return appliedmigrationversion, err
}

func (gateway *Gateway) ApplyMigrations(applyVersions uint) (int, error) {
	appliedmigrationversion, err := migrate.ExecMax(gateway.db, "postgres", *gateway.migration, migrate.Up, int(applyVersions))
	return appliedmigrationversion, err
}

func (gateway *Gateway) RollbackMigrations(applyVersions uint) (int, error) {
	appliedmigrationversion, err := migrate.ExecMax(gateway.db, "postgres", *gateway.migration, migrate.Down, int(applyVersions))
	return appliedmigrationversion, err
}

func main() {
	app := &cli.App{
		Name:        "personactl",
		Version:     "0.1.0",
		Description: "a cli application to moderate persona instance.",
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{"db"},
				Usage:   "Execute migration for database.",
				Subcommands: []*cli.Command{
					{
						Name:    "down",
						Aliases: []string{"d"},
						Usage:   "Roll back migrations",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:    "times",
								Aliases: []string{"t"},
								Value:   1,
							},
						},
						Action: func(context *cli.Context) error {
							var (
								appliedversion int
								err            error
							)
							configtree, err := config.LoadConfig(context.String("file"))
							if err != nil {
								return err
							}
							connectionData := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", configtree.Database.User, configtree.Database.Database, configtree.Database.Password, configtree.Database.SSL)
							db, err := sql.Open("postres", connectionData)
							if err != nil {
								return err
							}
							migrations := &migrate.FileMigrationSource{
								Dir: "../../migrations/postgres",
							}
							gateway := NewGateway(db, migrations)

							appliedversion, err = gateway.RollbackMigrations(context.Uint("times"))
							if err != nil {
								return err
							}
							fmt.Printf("applied %d migrations.", appliedversion)
							return nil
						},
					},
					{
						Name:    "up",
						Aliases: []string{"u"},
						Usage:   "",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:    "times",
								Aliases: []string{"t"},
								Value:   0,
								Usage:   "How many versions of migration to apply",
							},
							&cli.BoolFlag{
								Name:  "latest",
								Value: false,
								Usage: "migration to latest SQL schema",
							},
							&cli.StringFlag{
								Name:     "file",
								Aliases:  []string{"f"},
								FilePath: "../../config.toml",
							},
						},
						Action: func(context *cli.Context) error {
							var (
								appliedversion int
								err            error
							)
							configtree, err := config.LoadConfig(context.String("file"))
							if err != nil {
								return err
							}
							connectionData := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", configtree.Database.User, configtree.Database.Database, configtree.Database.Password, configtree.Database.SSL)
							db, err := sql.Open("postres", connectionData)
							if err != nil {
								return err
							}
							migrations := &migrate.FileMigrationSource{
								Dir: "../../migrations/postgres",
							}
							gateway := NewGateway(db, migrations)

							if context.Bool("latest") {
								appliedversion, err = gateway.ApplyMigrationsToLatest()
								if err != nil {
									return err
								}
							} else {
								appliedversion, err = gateway.ApplyMigrations(context.Uint("times"))
								if err != nil {
									return err
								}
							}
							fmt.Printf("applied %d migrations.", appliedversion)
							return nil
						},
					},
				},
			},
		},
	}
	_ = app.Run(os.Args)
}
