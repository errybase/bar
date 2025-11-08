package bar

import (
	"context"
	"reflect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type BelongsTo[T any] struct {
	Model        any
	RelationName string
}

func (r BelongsTo[T]) Get(ctx context.Context, db bun.IDB) (T, error) {
	var t T
	cols := r.appendRelModel(db, &t)
	err := db.NewSelect().Model(&t).WherePK(cols...).Scan(ctx)
	return t, err
}

func (r BelongsTo[T]) Set(ctx context.Context, db bun.IDB, t T) error {
	r.appendBaseModel(db, t)
	return Model(r.Model).Update(ctx, db, updateOmitZero)
}

func (r BelongsTo[T]) Create(ctx context.Context, db bun.IDB, t *T) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := Model(t).Create(ctx, tx); err != nil {
			return err
		}
		return r.Set(ctx, tx, *t)
	})
}

func (r BelongsTo[T]) appendRelModel(db bun.IDB, t *T) []string {
	rel := r.rel(db)
	bv := reflect.ValueOf(r.Model).Elem()
	tv := reflect.ValueOf(t).Elem()

	var cols []string
	for i, joinPK := range rel.JoinPKs {
		cols = append(cols, joinPK.Name)
		basePK := rel.BasePKs[i]
		tv.FieldByName(joinPK.GoName).Set(bv.FieldByName(basePK.GoName))
	}

	return cols
}

func (r BelongsTo[T]) appendBaseModel(db bun.IDB, t T) []string {
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

func (r BelongsTo[T]) rel(db bun.IDB) *schema.Relation {
	for name, rel := range r.baseTable(db).Relations {
		if name == r.RelationName {
			return rel
		}
	}
	return nil
}

func (r BelongsTo[T]) baseTable(db bun.IDB) *schema.Table {
	return db.Dialect().Tables().ByModel(r.baseModel())
}

func (r BelongsTo[T]) baseModel() string {
	return reflect.TypeOf(r.Model).Elem().Name()
}

func updateOmitZero(uq *bun.UpdateQuery) *bun.UpdateQuery {
	return uq.OmitZero()
}
