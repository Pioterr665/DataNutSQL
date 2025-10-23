package drivers

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PG struct {
	Conn *pgx.Conn
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
