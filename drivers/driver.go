package drivers

import (
	"context"
)

type SchemaInfo map[string][]string

type DatabaseDriver interface {
	Connect(connString string) error
	Close() error
	GetSchemaInfo(ctx context.Context) (SchemaInfo, error)
	QueryTable(ctx context.Context, query string, args ...any) ([][]any, error)
	QueryTableWithHeaders(ctx context.Context, query string, args ...any) ([][]any, error)
}
