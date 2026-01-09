package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/config"
)

func (ss *SensorsService) initMysqlConnection(ctx context.Context, opts config.DatabaseConnection) (err error) {
	_, found := ss.MysqlConnections[opts.Name]
	if found {
		return DuplicateConnectionError
	}
	db, err := sql.Open("mysql", opts.DatabaseConnectionString)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(opts.MaxOpenCons)
	db.SetMaxIdleConns(opts.MaxIdleCons)
	if opts.MaxOpenCons != opts.MaxIdleCons {
		log.Warn().Msgf("According to https://github.com/go-sql-driver/mysql?tab=readme-ov-file#important-settings"+
			" it is recommended to make `max_open_cons: %v` and `max_idle_cons: %v` equal for connection %s",
			opts.MaxOpenCons, opts.MaxIdleCons, opts.Name)
	}
	con, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	ss.MysqlConnections[opts.Name] = con
	log.Info().Msgf("MySQL/MariaDB database connection %s is established", opts.Name)
	return nil
}
