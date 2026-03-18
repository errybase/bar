package bar

import (
	"context"
	"database/sql"
)

func Model[T any](mod T) *model[T] {
	return &model[T]{mod}
}

type model[T any] struct {
	mod T
}

func (m *model[T]) Find(ctx context.Context) error {
	return DBFromContext(ctx).NewSelect().Model(m.mod).WherePK().Scan(ctx)
}

func (m *model[T]) Create(ctx context.Context) (err error) {
	err = m.validate(ctx)
	if err == nil {
		_, err = DBFromContext(ctx).NewInsert().Model(m.mod).Exec(ctx)
	}
	return
}

func (m *model[T]) FindOrCreate(ctx context.Context) error {
	if err := m.Find(ctx); err == sql.ErrNoRows {
		return m.Create(ctx)
	} else {
		return err
	}
}

func (m *model[T]) Update(ctx context.Context) error {
	if err := m.validate(ctx); err != nil {
		return err
	}

	r, err := DBFromContext(ctx).NewUpdate().Model(m.mod).WherePK().Exec(ctx)
	return handleResult(r, err)
}

func (m *model[T]) Save(ctx context.Context) error {
	if err := m.validate(ctx); err != nil {
		return err
	}

	db := DBFromContext(ctx)

	r, err := db.NewUpdate().Model(m.mod).WherePK().Exec(ctx)
	if err := handleResult(r, err); err == sql.ErrNoRows {
		_, err := db.NewInsert().Model(m.mod).Exec(ctx)
		return err
	} else {
		return err
	}
}

func (m *model[T]) Delete(ctx context.Context) error {
	r, err := DBFromContext(ctx).NewDelete().Model(m.mod).WherePK().Exec(ctx)
	return handleResult(r, err)
}

func (m *model[T]) validate(ctx context.Context) (err error) {
	if v := ValidateFromContext(ctx); v != nil {
		err = v.StructCtx(ctx, m.mod)
	}
	return
}

func handleResult(r sql.Result, err error) error {
	if err != nil {
		return err
	} else if i, err := r.RowsAffected(); err != nil {
		return err
	} else if i == 0 {
		return sql.ErrNoRows
	} else {
		return nil
	}
}
