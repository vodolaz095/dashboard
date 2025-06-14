package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/config"
)

func (ss *SensorsService) initPostgresConnection(ctx context.Context, opts config.DatabaseConnection) (err error) {
	_, found := ss.PostgresqlConnections[opts.Name]
	if found {
		return DuplicateConnectionError
	}
	db, err := sql.Open("pgx", opts.DatabaseConnectionString)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(opts.MaxOpenCons)
	db.SetMaxIdleConns(opts.MaxIdleCons)
	con, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	ss.PostgresqlConnections[opts.Name] = con
	log.Info().Msgf("Postgres database connection %s is established", opts.Name)
	return nil
}
