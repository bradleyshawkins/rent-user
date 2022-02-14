package postgres

import (
	"database/sql"
	"embed"
	"github.com/bradleyshawkins/berror"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log"

	_ "github.com/lib/pq"
)

var (
	//go:embed schema/*.sql
	migrationFiles embed.FS
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) (*Database, error) {
	database := &Database{db: db}
	err := database.migrate()
	if err != nil {
		return nil, err
	}
	return database, nil
}

func (d *Database) migrate() error {
	log.Println("Running database migrations...")
	defer log.Println("Completed database migrations")

	migrations, err := iofs.New(migrationFiles, "schema")
	if err != nil {
		return berror.WrapInternal(err, "unable to instantiate database migration type")
	}
	pdb, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		return berror.WrapInternal(err, "unable to instantiate database migration type")
	}

	m, err := migrate.NewWithInstance("iofs", migrations, "postgres", pdb)
	if err != nil {
		return berror.WrapInternal(err, "unable to instantiate migration type")
	}

	err = m.Up()
	if err != nil {
		return berror.WrapInternal(err, "unable to migrate database")
	}
	return nil
}

func (d *Database) begin() (*transaction, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	return &transaction{tx: tx}, nil
}

type transaction struct {
	tx *sql.Tx
}

func (t *transaction) commit() error {
	err := t.tx.Commit()
	if err != nil {
		return berror.WrapInternal(err, "unable to commit transaction")
	}
	return nil
}

func (t *transaction) rollback() error {
	err := t.tx.Rollback()
	if err != nil {
		return berror.WrapInternal(err, "unable to rollback transaction")
	}
	return nil
}
