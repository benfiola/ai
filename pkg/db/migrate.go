package db

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*
var migrationsFs embed.FS

type Version int

var VersionBase = Version(0)
var VersionHead = Version(-1)

func (db *DB) Migrate(version Version) error {
	source, err := iofs.New(migrationsFs, "")
	if err != nil {
		return err
	}

	var driver database.Driver
	switch db.URL.Scheme {
	case "sqlite":
		driver, err = sqlite3.WithInstance(db.Pool, &sqlite3.Config{})
	case "postgres":
		driver, err = pgx.WithInstance(db.Pool, &pgx.Config{})
	default:
		err = fmt.Errorf("unrecognized db type %s", db.URL.Scheme)
	}
	if err != nil {
		return err
	}

	instance, err := migrate.NewWithInstance("iofs", source, "ai", driver)
	if err != nil {
		return err
	}

	switch version {
	case VersionBase:
		return instance.Down()
	case VersionHead:
		return instance.Up()
	default:
		return instance.Migrate(uint(version))
	}
}
