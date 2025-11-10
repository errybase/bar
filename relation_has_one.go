package bar

import (
	"context"

	"github.com/uptrace/bun"
)

type HasOne[T any] relation[T]

func (r HasOne[T]) Get(ctx context.Context, db bun.IDB) (T, error) {
	var t T
	rel := relation[T](r).rel(db)
	rel.appendRelModel(&t)
	cols := rel.joinCols()
	err := db.NewSelect().Model(&t).WherePK(cols...).Scan(ctx)
	return t, err
}

func (r HasOne[T]) Set(ctx context.Context, db bun.IDB, t *T) error {
	relation[T](r).rel(db).appendRelModel(t)
	return Model(t).Update(ctx, db)
}

func (r HasOne[T]) Create(ctx context.Context, db bun.IDB, t *T) error {
	relation[T](r).rel(db).appendRelModel(t)
	return Model(t).Create(ctx, db)
}

func (r HasOne[T]) Update(ctx context.Context, db bun.IDB, t *T) error {
	cols := relation[T](r).rel(db).joinCols()
	_, err := db.NewUpdate().Model(t).WherePK(cols...).ExcludeColumn(cols...).Exec(ctx)
	return err
}

func (r HasOne[T]) Delete(ctx context.Context, db bun.IDB, t *T) error {
	cols := relation[T](r).rel(db).joinCols()
	_, err := db.NewDelete().Model(t).WherePK(cols...).Exec(ctx)
	return err
}
