package bar

import "github.com/uptrace/bun"

type Values map[string]any

func (values Values) QueryBuilder(qb bun.QueryBuilder) bun.QueryBuilder {
	for col, value := range values {
		qb = qb.Where("? = ?", bun.Ident(col), value)
	}
	return qb
}

func (values Values) InsertQuery(iq *bun.InsertQuery) *bun.InsertQuery {
	for col, value := range values {
		iq = iq.Value(col, "?", value)
	}
	return iq
}

func (values Values) UpdateQuery(uq *bun.UpdateQuery) *bun.UpdateQuery {
	for col, value := range values {
		uq = uq.Value(col, "?", value)
	}
	return uq
}
