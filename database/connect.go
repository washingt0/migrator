package database

import (
	"database/sql"
	"migrator/config"

	// driver used to connect with postgres
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

	if cfg.DestinationURI != "" {
		dstDB, err = sql.Open("postgres", cfg.DestinationURI)
		if err != nil {
			return
		}
	}

	return
}
