package bar_test

import (
	"database/sql"
	"testing"

	"github.com/errybase/bar"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type User struct {
	ID   int64 `bun:",pk,autoincrement"`
	Name string
}

func newDB(t *testing.T) (db *bun.DB) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	bar.MayPanic(err)
	db = bun.NewDB(sqldb, sqlitedialect.New()).
		WithQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	bar.MayPanic(db.ResetModel(t.Context(),
		(*User)(nil),
	))
	return
}

func TestCreate(t *testing.T) {
	db := newDB(t)

	john := &User{Name: "John"}
	if err := bar.Model(john).Create(t.Context(), db); err != nil {
		t.Error(err)
	}
}
