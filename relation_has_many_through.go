package bar

import (
	"context"
	"reflect"

	"github.com/uptrace/bun"
)

type HasManyThrough[R any, T any] relation[R]

func (r HasManyThrough[R, T]) All(ctx context.Context, db bun.IDB, fns ...func(*bun.SelectQuery) *bun.SelectQuery) (models []R, err error) {
	rel := relation[T](r).rel(db)
	var cols []string

	q := db.NewSelect().Table(rel.M2MTable.Name)
	for i, m2mJoinPK := range rel.M2MJoinPKs {
		joinPK := rel.JoinPKs[i]
		cols = append(cols, joinPK.Name)
		q = q.ColumnExpr("? AS ?", bun.Ident(m2mJoinPK.Name), bun.Ident(joinPK.Name))
	}
	for i, m2mBasePK := range rel.M2MBasePKs {
		basePK := rel.BasePKs[i]
		q = q.Where("? = ?", bun.Ident(m2mBasePK.Name), rel.v.FieldByName(basePK.GoName).Interface())
	}

	if e := q.Scan(ctx, &models); e != nil {
		err = e
	} else if len(models) > 0 {
		err = db.NewSelect().Model(&models).WherePK(cols...).Apply(fns...).Scan(ctx)
	}

	return
}

func (r HasManyThrough[R, T]) First(ctx context.Context, db bun.IDB) (m R, err error) {
	if models, e := r.All(ctx, db, func(sq *bun.SelectQuery) *bun.SelectQuery {
		return sq.Limit(1)
	}); e != nil {
		err = e
	} else {
		m = models[0]
	}
	return
}

func (r HasManyThrough[R, T]) Create(ctx context.Context, db bun.IDB, models ...*R) error {
	rel := relation[R](r).rel(db)

	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := Model(&models).Create(ctx, tx); err != nil {
			return err
		}

		var m2mModels []*T
		for _, m := range models {
			m2mModel := new(T)
			v := reflect.ValueOf(m2mModel).Elem()
			mv := reflect.ValueOf(m).Elem()

			for i, m2mBasePK := range rel.M2MBasePKs {
				basePK := rel.BasePKs[i]
				v.FieldByName(m2mBasePK.GoName).Set(rel.v.FieldByName(basePK.GoName))
			}

			for i, m2mJoinPK := range rel.M2MJoinPKs {
				joinPK := rel.JoinPKs[i]
				v.FieldByName(m2mJoinPK.GoName).Set(mv.FieldByName(joinPK.GoName))
			}

			m2mModels = append(m2mModels, m2mModel)
		}

		if err := Model(&m2mModels).Create(ctx, tx); err != nil {
			return err
		}

		return nil
	})
}
