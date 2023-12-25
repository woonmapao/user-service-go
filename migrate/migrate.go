package main

import (
	"log"

	"github.com/woonmapao/user-service-go/initializer"
	"github.com/woonmapao/user-service-go/models"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DBInitializer()
}

func main() {

	err := initializer.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to perform auto migration: &v", err)
	}
}
