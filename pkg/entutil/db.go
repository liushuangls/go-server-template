package entutil

import (
	"time"

	"entgo.io/ent/dialect/sql"
)

type Config struct {
	Dialect     string `koanf:"Dialect"`
	DSN         string `koanf:"DSN"`
	MaxIdle     int    `koanf:"MaxIdle"`
	MaxActive   int    `koanf:"MaxActive"`
	MaxLifetime int    `koanf:"MaxLifetime"`
}

func NewDriver(c *Config) (*sql.Driver, error) {
	drv, err := sql.Open(c.Dialect, c.DSN)
	if err != nil {
		return nil, err
	}
	db := drv.DB()
	db.SetMaxIdleConns(c.MaxIdle)
	db.SetMaxOpenConns(c.MaxActive)
	db.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return drv, nil
}
