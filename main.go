package main

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/user-service-go/initializer"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DBInitializer()
}

func main() {

	r := gin.Default()
	r.Run()

}
