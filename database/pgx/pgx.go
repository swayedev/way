package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/swayedev/way/database/config"
)

func Connect(uri string) (*pgx.Conn, error) {
	if uri == "" && config.GetDbUri() != "" {
		uri = config.GetDbUri()
	}

	if uri == "" {
		uri = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.GetDbUser(), config.GetDbPassword(), config.GetDbHost(), config.GetDbPort(), config.GetDbName())
	}

	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
