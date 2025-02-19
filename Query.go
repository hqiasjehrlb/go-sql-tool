package sqltool

import (
	"context"
)

type Query[T any] struct {
	DBConn
}

func (db *Query[T]) QueryRow(ctx context.Context, query string, args ...any) (*T, ExtFieldMap, error) {
	return QueryRow[T](ctx, db, query, args...)
}

func (db *Query[T]) QueryRows(ctx context.Context, query string, args ...any) ([]T, []ExtFieldMap, error) {
	return QueryRows[T](ctx, db, query, args...)
}
