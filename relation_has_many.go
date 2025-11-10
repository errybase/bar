package bar

import (
	"context"

	"github.com/uptrace/bun"
)

type HasMany[T any] relation[T]

func (r HasMany[T]) All(ctx context.Context, db bun.IDB, fns ...func(*bun.SelectQuery) *bun.SelectQuery) (models []T, err error) {
	rel := relation[T](r).rel(db)
	if rel.M2MTable == nil {
		err = db.NewSelect().Model(&models).Where("(?) IN (?)", bun.In(rel.joinIdentCols()), bun.In(rel.basePKValues())).Apply(fns...).Scan(ctx)
	} else {
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
		} else {
			err = db.NewSelect().Model(&models).WherePK(cols...).Scan(ctx)
		}
	}
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
