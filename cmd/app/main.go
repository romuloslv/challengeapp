package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/romuloslv/go-rest-api/cmd/app/config"
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

	// Register our service handlers to the router
	router := gin.Default()
	accountService.RegisterHandlers(router)

	// Start the server
	router.Run()
}
