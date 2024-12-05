package drivers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectWithSQLX(opt *DriverOptions) (*sqlx.DB, error) {
	conn, err := sqlx.Open(opt.DriverName, opt.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("driver.Connect - sqlx.Open: %v", err)
	}

	conn.SetMaxOpenConns(opt.MaxOpenConnections)
	conn.SetMaxIdleConns(opt.MaxIdleConnections)
	conn.SetConnMaxLifetime(opt.ConnMaxLifetime)

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("driver.Connect - conn.Ping: %v", err)
	}

	return conn, nil
}
