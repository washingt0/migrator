package main

import (
	"log"
	"os"

	"migrator/config"
	"migrator/database"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var (
		outFile *os.File
		err     error
		tables  *database.Tables
	)

	config.LoadConfig()

	database.Connect()
	tables, err = database.Analize()
	if err != nil {
		log.Fatal(err)
	}

	if config.GetConfig().DestinationFile != "" {
		outFile, err = os.OpenFile(config.GetConfig().DestinationFile, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := database.Migrate(tables, outFile); err != nil {
		log.Fatal(err)
	}
}
