package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	_ "github.com/lib/pq"
)

var db DB

type SQLOperations interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ValidForPostgres() bool
}

type pgSQLOperations struct {
	*sql.Tx
}

func (o *pgSQLOperations) ValidForPostgres() bool {
	return true
}

type DB interface {
	SQLOperations
	Begin() (*sql.Tx, error)
	Close() error
	Ping() error
	InTransaction(ctx context.Context, operations func(context.Context, SQLOperations) error) error
	Valid() bool
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

type appDB struct {
	*sql.DB
	valid bool
}

func (db *appDB) InTransaction(ctx context.Context, operations func(context.Context, SQLOperations) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	sqlOperations := &pgSQLOperations{
		Tx: tx,
	}

	if err = operations(ctx, sqlOperations); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	return tx.Commit()
}

func (db *appDB) ValidForPostgres() bool {
	return true
}

func (db *appDB) Valid() bool {
	return db.valid
}

func InitDB(databaseURL string) DB {
	return InitDBWithURL(
		databaseURL,
	)
}

func InitDBWithURL(databaseURL string) DB {
	pgDB := newPostgresDBWithURL(databaseURL)

	db = &appDB{
		DB:    pgDB,
		valid: true,
	}

	err := db.Ping()
	if err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	return db
}

func GetDB() DB {
	return db
}

func newPostgresDBWithURL(databaseURL string) *sql.DB {
	if databaseURL == "" {
		logger.Fatal("database url is empty and is required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		logger.Fatalf("sql.Open function call failed: [%v]", err)
	}

	// Set connection pool configurations
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	// Verify the database connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Database connection failed: [%v]", err)
	}

	return db
}
