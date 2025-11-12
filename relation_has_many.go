package bar

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type HasMany[T any] relation[T]

func (r HasMany[T]) All(ctx context.Context, db bun.IDB, fns ...func(*bun.SelectQuery) *bun.SelectQuery) (models []T, err error) {
	rel := relation[T](r).rel(db)
	err = db.NewSelect().Model(&models).Where("(?) IN (?)", bun.In(rel.joinIdentCols()), bun.In(rel.basePKValues())).Apply(fns...).Scan(ctx)
	return
}

func (r HasMany[T]) First(ctx context.Context, db bun.IDB, fns ...func(*bun.SelectQuery) *bun.SelectQuery) (m T, err error) {
	fns = append(fns, queryLimit1)
	if models, e := r.All(ctx, db, fns...); e != nil {
		err = e
	} else if len(models) > 0 {
		m = models[0]
	} else {
		err = sql.ErrNoRows
	}
	return
}

func (r HasMany[T]) Create(ctx context.Context, db bun.IDB, models ...*T) error {
	rel := relation[T](r).rel(db)
	for _, m := range models {
		rel.appendRelModel(m)
	}
	return Model(&models).Create(ctx, db)
}

func queryLimit1(sq *bun.SelectQuery) *bun.SelectQuery {
	return sq.Limit(1)
}
