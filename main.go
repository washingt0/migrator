package main

import (
	"log"

	"migrator/config"
	"migrator/database"

	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.LoadConfig()

	database.Connect()
	tables, err := database.Analize()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(tables); err != nil {
		log.Fatal(err)
	}
}
