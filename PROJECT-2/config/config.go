package config

import (
	"PROJECT-2/database"
	"database/sql"
	"log"
)

type Config struct {
	DB *sql.DB
}

func NewConfig() *Config {
	log.Println("Access creating config")
	return &Config{
		database.NewDB(),
	}
}
