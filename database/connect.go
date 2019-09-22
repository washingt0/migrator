package database

import (
	"database/sql"
	"migrator/config"

	_ "github.com/lib/pq"
)

var (
	srcDB *sql.DB
	dstDB *sql.DB
)

// Connect TODO:
func Connect() (err error) {
	cfg := config.GetConfig()
	srcDB, err = sql.Open("postgres", cfg.SourceURI)
	if err != nil {
		return
	}

	dstDB, err = sql.Open("postgres", cfg.DestinationURI)
	if err != nil {
		return
	}

	return
}
