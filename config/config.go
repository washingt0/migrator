package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var conf *Config

// Config holds all needed data to perform migration
type Config struct {
	SourceURI       string   `json:"source_uri"`
	DestinationURI  string   `json:"destination_uri"`
	DestinationFile string   `json:"destination_file"`
	Tables          []string `json:"tables"`
	RecordLimit     int      `json:"record_limit"`
}

// LoadConfig read and parse configuration file
func LoadConfig() {
	path := "config.json"

	if val, set := os.LookupEnv(`MIGRATOR_CONFIG`); set && val != "" {
		path = val
	} else {
		log.Println("MIGRATOR_CONFIG enviroment variable is not set or empty, using default: ", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	if cfg.SourceURI == "" {
		log.Fatal("Source URI is not specified")
	}

	if cfg.DestinationURI == "" && cfg.DestinationFile == "" {
		log.Fatal("Destination URI is not specified")
	}

	if len(cfg.Tables) == 0 {
		log.Fatal("Tables is empty")
	}

	conf = cfg
}

// GetConfig returns a copy of loaded config
func GetConfig() Config {
	if conf == nil {
		log.Fatal("Settings not loaded")
	}
	return *conf
}
