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

type relation[T any] struct {
	Model        any
	RelationName string
}

func (r relation[T]) appendRelModel(db bun.IDB, t *T) {
	rel := r.rel(db)
	bv := reflect.ValueOf(r.Model).Elem()
	tv := reflect.ValueOf(t).Elem()

	for i, joinPK := range rel.JoinPKs {
		basePK := rel.BasePKs[i]
		tv.FieldByName(joinPK.GoName).Set(bv.FieldByName(basePK.GoName))
	}
}

func (r relation[T]) appendBaseModel(db bun.IDB, t T) []string {
	rel := r.rel(db)
	bv := reflect.ValueOf(r.Model).Elem()
	tv := reflect.ValueOf(t)

	var cols []string
	for i, basePK := range rel.BasePKs {
		cols = append(cols, basePK.Name)
		joinPK := rel.JoinPKs[i]
		bv.FieldByName(basePK.GoName).Set(tv.FieldByName(joinPK.GoName))
	}

	return cols
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

func (r relation[T]) baseTable(db bun.IDB) *schema.Table {
	return db.Dialect().Tables().ByModel(r.baseModel())
}

func (r relation[T]) baseModel() string {
	return reflect.TypeOf(r.Model).Elem().Name()
}
