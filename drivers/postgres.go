package drivers

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PG struct {
	Conn *pgx.Conn
}

func (pg *PG) Connect(connString string) error {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return err
	}
	pg.Conn = conn
	return nil
}

func (pg *PG) Close() error {
	return pg.Conn.Close(context.Background())
}

func (pg *PG) GetSchemaInfo(ctx context.Context) (SchemaInfo, error) {
	return GetSchemaInfo(ctx, pg.Conn)
}

func (pg *PG) QueryTable(ctx context.Context, query string, args ...any) ([][]any, error) {
	return QueryTable(ctx, pg.Conn, query, args...)
}

func (pg *PG) QueryTableWithHeaders(ctx context.Context, query string, args ...any) ([][]any, error) {
	return QueryTableWithHeaders(ctx, pg.Conn, query, args...)
}

func GetSchemaInfo(ctx context.Context, conn *pgx.Conn) (SchemaInfo, error) {
	schema := make(SchemaInfo)
	rows, err := conn.Query(ctx, `
		SELECT table_name, column_name
		FROM information_schema.columns
		WHERE table_schema = 'public'
		ORDER BY table_name, ordinal_position
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table, column string
		if err := rows.Scan(&table, &column); err != nil {
			return nil, err
		}
		schema[table] = append(schema[table], column)
	}
	return schema, nil
}

func ConnectPG(conn_string string) (*PG, error) {
	conn, err := pgx.Connect(context.Background(), conn_string)
	if err != nil {
		return nil, err
	}
	return &PG{Conn: conn}, nil
}

func (pg *PG) ClosePG() error {
	return pg.Conn.Close(context.Background())
}

func QueryTable(ctx context.Context, conn *pgx.Conn, query string, args ...any) ([][]any, error) {
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results [][]any
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		results = append(results, values)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return results, nil
}

func QueryTableWithHeaders(ctx context.Context, conn *pgx.Conn, query string, args ...any) ([][]any, error) {
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fieldDescs := rows.FieldDescriptions()
	headers := make([]any, len(fieldDescs))
	for i, fd := range fieldDescs {
		headers[i] = string(fd.Name)
	}

	var results [][]any
	results = append(results, headers)
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		results = append(results, values)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return results, nil
}
