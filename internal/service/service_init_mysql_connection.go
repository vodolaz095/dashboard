package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/config"
)

func (ss *SensorsService) initMysqlConnection(ctx context.Context, opts config.DatabaseConnection) (err error) {
	_, found := ss.MysqlConnections[opts.Name]
	if found {
		return ErrDuplicateConnection
	}
	db, err := sql.Open("mysql", opts.DatabaseConnectionString)
	if err != nil {
		return fmt.Errorf("error opening mysql connection: %w", err)
	}
	db.SetMaxOpenConns(opts.MaxOpenCons)
	if opts.MaxIdleCons > 0 {
		db.SetMaxIdleConns(opts.MaxIdleCons)
	} else {
		db.SetMaxIdleConns(opts.MaxOpenCons)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("error pinging mysql: %w", err)
	}
	ss.MysqlConnections[opts.Name] = db
	log.Info().Msgf("MySQL/MariaDB database connection pool %s is established with %d max connections",
		opts.Name, opts.MaxOpenCons)
	return nil
}
