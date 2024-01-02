package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/swayedev/way/database/config"
)

func Connect() (*pgx.Conn, error) {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.GetDbUser(), config.GetDbPassword(), config.GetDbHost(), config.GetDbPort(), config.GetDbName())
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
