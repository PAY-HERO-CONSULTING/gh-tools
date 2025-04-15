package testutils

import (
	"context"
	"database/sql"
	"testing"

	"github.com/PAY-HERO-CONSULTING/gh-tools/db"
	"github.com/stretchr/testify/assert"
)

func WithTestDBs(
	ctx context.Context,
	testDB db.DB,
	f func(context.Context, db.DB, *assert.Assertions),
) func(t *testing.T) {
	return func(t *testing.T) {
		assert := assert.New(t)

		var dbTx *sql.Tx

		dB := db.NewTestDB(dbTx)

		if testDB != nil && testDB.Valid() {

			_, err := testDB.ExecContext(ctx, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
			assert.NoError(err)

			dbTx, err = testDB.Begin()
			assert.NoError(err)

			dB = db.NewTestDB(dbTx)
		}

		t.Cleanup(func() {
			if dbTx != nil {
				err := dbTx.Rollback()
				assert.NoError(err)
			}
		})

		f(ctx, dB, assert)
	}
}
