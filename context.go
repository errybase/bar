package bar

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type validateCtxKey struct{}

func ContextWithValidate(ctx context.Context, v *validator.Validate) context.Context {
	return context.WithValue(ctx, validateCtxKey{}, v)
}

func ValidateFromContext(ctx context.Context) *validator.Validate {
	if v := ctx.Value(validateCtxKey{}); v != nil {
		return v.(*validator.Validate)
	}
	return nil
}

type dbCtxKey struct{}

func ContextWithDB(ctx context.Context, db bun.IDB) context.Context {
	return context.WithValue(ctx, dbCtxKey{}, db)
}

func DBFromContext(ctx context.Context) bun.IDB {
	return ctx.Value(dbCtxKey{}).(bun.IDB)
}
