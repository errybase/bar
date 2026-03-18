package bar

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

func ModelFor[T any]() Model[T] {
	return Model[T]{}
}

type Model[T any] struct {
	parent *Model[T]
	values Values
	sq     func(*bun.SelectQuery) *bun.SelectQuery
}

func (m Model[T]) Where(values Values) Model[T] {
	return Model[T]{
		parent: &m,
		values: values,
	}
}

func (m Model[T]) Limit(n int) Model[T] {
	return Model[T]{
		parent: &m,
		sq: func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Limit(n)
		},
	}
}

func (m Model[T]) All(ctx context.Context) (models []T, err error) {
	err = DBFromContext(ctx).NewSelect().Model(&models).ApplyQueryBuilder(m.queryBuilder).Apply(m.selectQuery).Scan(ctx)
	return
}

func (m Model[T]) First(ctx context.Context) (*T, error) {
	if models, err := m.Limit(1).All(ctx); err != nil {
		return nil, err
	} else if len(models) == 0 {
		return nil, sql.ErrNoRows
	} else {
		return &models[0], nil
	}
}

func (m Model[T]) Create(ctx context.Context, values ...Values) (*T, error) {
	queries := []func(*bun.InsertQuery) *bun.InsertQuery{m.insertQuery}
	for _, v := range values {
		queries = append(queries, v.InsertQuery)
	}

	var model T
	_, err := DBFromContext(ctx).NewInsert().Model(&model).Apply(queries...).Exec(ctx)
	return &model, err
}

func (m Model[T]) Update(ctx context.Context, values ...Values) (*T, error) {
	var queries []func(*bun.UpdateQuery) *bun.UpdateQuery
	for _, v := range values {
		queries = append(queries, v.UpdateQuery)
	}

	var model T
	_, err := DBFromContext(ctx).NewUpdate().Model(&model).ApplyQueryBuilder(m.queryBuilder).Apply(queries...).Exec(ctx)
	return &model, err
}

func (m Model[T]) Delete(ctx context.Context) error {
	_, err := DBFromContext(ctx).NewDelete().Model((*T)(nil)).ApplyQueryBuilder(m.queryBuilder).Exec(ctx)
	return err
}

func (m Model[T]) queryBuilder(qb bun.QueryBuilder) bun.QueryBuilder {
	if p := m.parent; p != nil {
		qb = p.queryBuilder(qb)
	}
	if values := m.values; values != nil {
		qb = values.QueryBuilder(qb)
	}
	return qb
}

func (m Model[T]) selectQuery(sq *bun.SelectQuery) *bun.SelectQuery {
	if p := m.parent; p != nil {
		sq = p.selectQuery(sq)
	}
	if msq := m.sq; msq != nil {
		sq = msq(sq)
	}
	return sq
}

func (m Model[T]) insertQuery(iq *bun.InsertQuery) *bun.InsertQuery {
	if p := m.parent; p != nil {
		iq = p.insertQuery(iq)
	}
	if values := m.values; values != nil {
		iq = values.InsertQuery(iq)
	}
	return iq
}
