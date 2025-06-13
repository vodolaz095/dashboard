package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) initMysqlConnection(ctx context.Context, name, dsn string) (err error) {
	_, found := ss.MysqlConnections[name]
	if found {
		return DuplicateConnectionError
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(3)
	con, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	ss.MysqlConnections[name] = con
	log.Info().Msgf("MySQL/MariaDB database connection %s is established", name)
	return nil
}
