package main

import (
	"log"

	"github.com/romuloslv/go-rest-api/cmd/core/config"
	"github.com/romuloslv/go-rest-api/local/database"
)

func main() {
	// Read configuration
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Instantiates the database
	postgres, err := database.NewPostgres(cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password)
	if err != nil {
		log.Fatal(err.Error())
	}
}
