package bar

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

var _ bun.BeforeAppendModelHook = (*ModelWithTimestamps)(nil)

type ModelWithTimestamps struct {
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *ModelWithTimestamps) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

type BaseModel struct {
	ModelWithTimestamps

	ID int64 `bun:"id,pk,autoincrement"`
}

type modelQuery struct {
	model any
}

func Model(model any) modelQuery {
	return modelQuery{model}
}

func (q modelQuery) Create(ctx context.Context, db bun.IDB) error {
	_, err := db.NewInsert().Model(q.model).Exec(ctx)
	return err
}

func (q modelQuery) Update(ctx context.Context, db bun.IDB, fns ...func(*bun.UpdateQuery) *bun.UpdateQuery) error {
	_, err := db.NewUpdate().Model(q.model).WherePK().Apply(fns...).Exec(ctx)
	return err
}
