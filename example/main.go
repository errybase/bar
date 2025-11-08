package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/errybase/bar"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	sqldb, _ = sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	db       = bun.NewDB(sqldb, sqlitedialect.New())
	ctx      = context.Background()
)

func init() {
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	db.ResetModel(ctx,
		(*User)(nil),
	)
}

var (
	users = []User{
		{
			Name:   "admin",
			Emails: []string{"admin1@admin", "admin2@admin"},
		},
		{
			Name:   "root",
			Emails: []string{"root1@root", "root2@root"},
		},
	}
)

func main() {
	if err := bar.Model(&users).Create(ctx, db); err != nil {
		panic(err)
	} else {
		fmt.Println("created users:", users)
	}
}

type User struct {
	bar.BaseModel

	Name   string
	Emails []string
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v %v %v>", u.ID, u.Name, u.Emails, u.CreatedAt, u.UpdatedAt)
}
