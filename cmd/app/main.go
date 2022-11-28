package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/romuloslv/challengeapp/api/accounts"
	"github.com/romuloslv/challengeapp/cmd/app/config"
	_ "github.com/romuloslv/challengeapp/docs"
	"github.com/romuloslv/challengeapp/internal/database"
)

// @Title Go REST API
// @Version 1.0.0
// @Description This is a sample server for a Go REST API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
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

	// Instantiates the account service
	queries := database.New(postgres.DB)
	accountService := accounts.NewService(queries)

	// Register our service handlers to the router
	router := gin.Default()
	accountService.RegisterHandlers(router)

	// Start the server
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
