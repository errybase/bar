package bar

import (
	"context"

	"github.com/uptrace/bun"
)

type HasMany[T any] relation[T]

func (r HasMany[T]) All(ctx context.Context, db bun.IDB, fns ...func(*bun.SelectQuery) *bun.SelectQuery) (models []T, err error) {
	rel := relation[T](r).rel(db)
	err = db.NewSelect().Model(&models).Where("(?) IN (?)", bun.In(rel.joinIdentCols()), bun.In(rel.basePKValues())).Apply(fns...).Scan(ctx)
	return
}

func (r HasMany[T]) First(ctx context.Context, db bun.IDB) (m T, err error) {
	if models, e := r.All(ctx, db, func(sq *bun.SelectQuery) *bun.SelectQuery {
		return sq.Limit(1)
	}); e != nil {
		err = e
	} else {
		m = models[0]
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
