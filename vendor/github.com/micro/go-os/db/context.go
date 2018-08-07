package db

import (
	"golang.org/x/net/context"
)

type dbKey struct{}

func FromContext(ctx context.Context) (DB, bool) {
	c, ok := ctx.Value(dbKey{}).(DB)
	return c, ok
}

func NewContext(ctx context.Context, c DB) context.Context {
	return context.WithValue(ctx, dbKey{}, c)
}
