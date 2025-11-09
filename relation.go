package bar

import (
	"context"
	"reflect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type BelongsTo[T any] relation[T]

func (r BelongsTo[T]) Get(ctx context.Context, db bun.IDB) (T, error) {
	var t T
	relation[T](r).appendRelModel(db, &t)
	cols := relation[T](r).joinCols(db)
	err := db.NewSelect().Model(&t).WherePK(cols...).Scan(ctx)
	return t, err
}

func (r BelongsTo[T]) Set(ctx context.Context, db bun.IDB, t T) error {
	relation[T](r).appendBaseModel(db, t)
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

type HasOne[T any] relation[T]

func (r HasOne[T]) Get(ctx context.Context, db bun.IDB) (T, error) {
	var t T
	relation[T](r).appendRelModel(db, &t)
	cols := relation[T](r).joinCols(db)
	err := db.NewSelect().Model(&t).WherePK(cols...).Scan(ctx)
	return t, err
}

func (r HasOne[T]) Set(ctx context.Context, db bun.IDB, t *T) error {
	relation[T](r).appendRelModel(db, t)
	return Model(t).Update(ctx, db)
}

func (r HasOne[T]) Create(ctx context.Context, db bun.IDB, t *T) error {
	relation[T](r).appendRelModel(db, t)
	return Model(t).Create(ctx, db)
}

func (r HasOne[T]) Update(ctx context.Context, db bun.IDB, t *T) error {
	cols := relation[T](r).joinCols(db)
	_, err := db.NewUpdate().Model(t).WherePK(cols...).ExcludeColumn(cols...).Exec(ctx)
	return err
}

func (r HasOne[T]) Delete(ctx context.Context, db bun.IDB, t *T) error {
	cols := relation[T](r).joinCols(db)
	_, err := db.NewDelete().Model(t).WherePK(cols...).Exec(ctx)
	return err
}

type HasMany[T any] relation[T]

func (r HasMany[T]) All(ctx context.Context, db bun.IDB) ([]T, error) {
	var models []T
	err := db.NewSelect().Model(&models).ApplyQueryBuilder(relation[T](r).queryBuilder(db)).Scan(ctx)
	return models, err
}

func (r HasMany[T]) First(ctx context.Context, db bun.IDB) (T, error) {
	var model T
	err := db.NewSelect().Model(&model).ApplyQueryBuilder(relation[T](r).queryBuilder(db)).Scan(ctx)
	return model, err
}

type relation[T any] struct {
	Model        any
	RelationName string
}

func (r relation[T]) queryBuilder(db bun.IDB) func(bun.QueryBuilder) bun.QueryBuilder {
	joinTable := r.rel(db).JoinTable
	tableName := joinTable.Name
	if joinTable.Alias != "" {
		tableName = joinTable.Alias
	}

	return func(qb bun.QueryBuilder) bun.QueryBuilder {
		q := qb
		for _, f := range r.fields(db) {
			q = q.Where("?.? = ?", bun.Ident(tableName), bun.Ident(f.join.Name), f.value.Interface())
		}
		return q
	}
}

func (r relation[T]) appendRelModel(db bun.IDB, t *T) {
	v := reflect.ValueOf(t).Elem()
	for _, f := range r.fields(db) {
		v.FieldByName(f.join.GoName).Set(f.value)
	}
}

func (r relation[T]) appendBaseModel(db bun.IDB, t T) {
	v := reflect.ValueOf(t)
	for _, f := range r.fields(db) {
		r.baseValue().FieldByName(f.base.GoName).Set(v.FieldByName(f.join.GoName))
	}
}

func (r relation[T]) fields(db bun.IDB) []field {
	fields, rel, baseValue := []field{}, r.rel(db), r.baseValue()
	for i, joinPK := range rel.JoinPKs {
		basePK := rel.BasePKs[i]
		fields = append(fields, field{
			base:  basePK,
			join:  joinPK,
			value: baseValue.FieldByName(basePK.GoName),
		})
	}
	return fields
}

func (r relation[T]) joinCols(db bun.IDB) []string {
	var cols []string
	for _, pk := range r.rel(db).JoinPKs {
		cols = append(cols, pk.Name)
	}
	return cols
}

func (r relation[T]) rel(db bun.IDB) *schema.Relation {
	for name, rel := range r.baseTable(db).Relations {
		if name == r.RelationName {
			return rel
		}
	}
	return nil
}

func (r relation[T]) baseValue() reflect.Value {
	return reflect.ValueOf(r.Model).Elem()
}

func (r relation[T]) baseTable(db bun.IDB) *schema.Table {
	return db.Dialect().Tables().ByModel(r.baseModel())
}

func (r relation[T]) baseModel() string {
	return reflect.TypeOf(r.Model).Elem().Name()
}

type field struct {
	base  *schema.Field
	join  *schema.Field
	value reflect.Value
}
