package bar

import (
	"context"

	"github.com/uptrace/bun"
)

type BelongsTo[T any] relation[T]

func (r BelongsTo[T]) Get(ctx context.Context, db bun.IDB) (T, error) {
	var t T
	rel := relation[T](r).rel(db)
	rel.appendRelModel(&t)
	cols := rel.joinCols()
	err := db.NewSelect().Model(&t).WherePK(cols...).Scan(ctx)
	return t, err
}

func (r BelongsTo[T]) Set(ctx context.Context, db bun.IDB, t T) error {
	relation[T](r).rel(db).appendBaseModel(t)
	return Model(r.Model).Update(ctx, db)
}

func (r BelongsTo[T]) Create(ctx context.Context, db bun.IDB, t *T) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := Model(t).Create(ctx, tx); err != nil {
			return err
		}
		return r.Set(ctx, tx, *t)
	})
}
