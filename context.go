package bar

import (
	"context"

	"github.com/go-playground/validator/v10"
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
