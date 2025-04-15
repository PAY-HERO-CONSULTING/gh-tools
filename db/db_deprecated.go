package db

// import (
// 	"context"
// 	"database/sql"

// 	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
// 	_ "github.com/lib/pq"
// 	// _ "github.com/jackc/pgx/v5" // TODO: Migrate to PGX
// )

// var db DB

// type SQLOperations interface {
// 	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
// 	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
// 	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
// 	ValidForPostgres() bool
// }

// type pgSQLOperations struct {
// 	*sql.Tx
// }

// func (o *pgSQLOperations) ValidForPostgres() bool {
// 	return true
// }

// type DB interface {
// 	SQLOperations
// 	Begin() (*sql.Tx, error)
// 	Close() error
// 	Ping() error
// 	InTransaction(ctx context.Context, operations func(context.Context, SQLOperations) error) error
// 	Valid() bool
// }

// type RowScanner interface {
// 	Scan(dest ...interface{}) error
// }

// type appDB struct {
// 	*sql.DB
// 	valid bool
// }

// func (db *appDB) InTransaction(ctx context.Context, operations func(context.Context, SQLOperations) error) error {

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	sqlOperations := &pgSQLOperations{
// 		Tx: tx,
// 	}

// 	if err = operations(ctx, sqlOperations); err != nil {
// 		if rollbackErr := tx.Rollback(); rollbackErr != nil {
// 			return rollbackErr
// 		}

// 		return err
// 	}

// 	return tx.Commit()
// }

// func (db *appDB) ValidForPostgres() bool {
// 	return true
// }

// func (db *appDB) Valid() bool {
// 	return db.valid
// }

// func InitDB(databaseURL string) DB {
// 	return InitDBWithURL(
// 		databaseURL,
// 	)
// }

// var STAGE_DATABASE_URL string

// func StageInitDB() DB {
// 	return InitDBWithURL(STAGE_DATABASE_URL)
// }

// func InitDBWithURL(databaseURL string) DB {

// 	dB := newPGDBWithURL(databaseURL)

// 	db = &appDB{
// 		DB:    dB,
// 		valid: true,
// 	}

// 	err := db.Ping()
// 	if err != nil {
// 		logger.Fatalf("db ping failed: %v", err)
// 	}

// 	return db
// }

// func GetDB() DB {
// 	return db
// }

// func newPGDBWithURL(databaseURL string) *sql.DB {
// 	if databaseURL == "" {
// 		logger.Fatal("database url is empty and is required")
// 	}

// 	db, err := sql.Open("postgres", databaseURL)
// 	if err != nil {
// 		logger.Fatalf("sql open fxn call failed: [%v]", err)
// 	}

// 	return db
// }
