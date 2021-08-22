package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"

	// register the postgres database/sql driver
	_ "github.com/lib/pq"
)

type Config struct {
	Host       string
	DBName     string
	User       string
	Pass       string
	DisableSSL bool
}

// Open knows how to open database connection
func Open(cfg Config) (*sqlx.DB, error) {
	q := url.Values{}

	if cfg.DisableSSL {
		q.Set("sslmode", "disable")
	}

	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Pass),
		Host:     cfg.Host,
		Path:     cfg.DBName,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
