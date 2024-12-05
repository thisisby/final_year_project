package drivers

import "time"

type DriverOptions struct {
	DriverName         string
	DataSourceName     string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
}

func NewDriverOptions(
	driverName,
	dataSourceName string,
	maxOpenConnections,
	maxIdleConnections int,
	connMaxLifetime time.Duration,
) *DriverOptions {
	return &DriverOptions{
		DriverName:         driverName,
		DataSourceName:     dataSourceName,
		MaxOpenConnections: maxOpenConnections,
		MaxIdleConnections: maxIdleConnections,
		ConnMaxLifetime:    connMaxLifetime,
	}
}
