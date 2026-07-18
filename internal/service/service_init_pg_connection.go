package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/config"
)

func (ss *SensorsService) initPostgresConnection(ctx context.Context, opts config.DatabaseConnection) (err error) {
	_, found := ss.PostgresqlConnections[opts.Name]
	if found {
		return ErrDuplicateConnection
	}
	db, err := sql.Open("pgx", opts.DatabaseConnectionString)
	if err != nil {
		return fmt.Errorf("error opening postgresql connection: %w", err)
	}
	db.SetMaxOpenConns(opts.MaxOpenCons)
	if opts.MaxIdleCons > 0 {
		if opts.MaxIdleCons > opts.MaxOpenCons {
			opts.MaxIdleCons = opts.MaxOpenCons
		}
		db.SetMaxIdleConns(opts.MaxIdleCons)
	} else {
		idle := opts.MaxOpenCons / 4
		if idle < 2 {
			idle = 2
		}
		db.SetMaxIdleConns(idle)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("error pinging postgresql: %w", err)
	}
	ss.PostgresqlConnections[opts.Name] = db
	log.Info().Msgf("PostgreSQL database connection pool %s is established with %d max connections and %d idle connections",
		opts.Name, opts.MaxOpenCons, opts.MaxIdleCons)
	return nil
}
