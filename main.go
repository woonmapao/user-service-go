package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/user-service-go/initializer"
	"github.com/woonmapao/user-service-go/routes"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DBInitializer()
}

func main() {

	r := gin.Default()

	// Setup routes
	routes.SetupUserRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}
