package bar

import (
	"context"

	"github.com/uptrace/bun"
)

func Model[T any](mod T) *model[T] {
	return &model[T]{mod}
}

type model[T any] struct {
	mod T
}

func (m *model[T]) Create(ctx context.Context, db bun.IDB) (err error) {
	err = m.validate(ctx)
	if err == nil {
		_, err = db.NewInsert().Model(m.mod).Exec(ctx)
	}
	return
}

func (m *model[T]) validate(ctx context.Context) (err error) {
	if v := ValidateFromContext(ctx); v != nil {
		err = v.StructCtx(ctx, m.mod)
	}
	return
}
