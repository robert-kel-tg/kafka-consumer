package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/robertke/kafka-consumer/pkg/infrastructure/config"
)

func Connect(driverName string, dbConf *config.DBConfig) (*sqlx.DB, error) {
	conn, err := sqlx.Connect(driverName, dbConf.ConnectionURL)
	if err != nil {
		return nil, err
	}

	// Set:
	// ConnMaxLifetime
	// MaxIdleConn
	// MaxOpenConn

	return conn, nil
}
