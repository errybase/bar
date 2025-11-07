package bar

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

var _ bun.BeforeUpdateHook = (*ModelWithTimestamps)(nil)

type ModelWithTimestamps struct {
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *ModelWithTimestamps) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	m.UpdatedAt = time.Now()
	return nil
}

type BaseModel struct {
	ModelWithTimestamps

	ID int64 `bun:"id,pk,autoincrement"`
}
