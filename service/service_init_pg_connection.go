package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) initPostgresConnection(ctx context.Context, name, dsn string) (err error) {
	_, found := ss.PostgresqlConnections[name]
	if found {
		return DuplicateConnectionError
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	con, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	ss.PostgresqlConnections[name] = con
	log.Info().Msgf("Postgres database connection %s is established", name)
	return nil
}
