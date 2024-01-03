package routes

import (
	"github.com/gin-gonic/gin"
	h "github.com/woonmapao/user-service-go/handlers"
)

func SetupUserRoutes(router *gin.Engine) {

	userGroup := router.Group("/users")
	{
		userGroup.GET("/", h.GetUsersHandler)
		userGroup.POST("/", h.AddUserHandler)

		userGroup.GET("/:id", h.GetUserHandler)
		userGroup.PUT("/:id", h.UpdateUserHandler)
		userGroup.DELETE("/:id", h.DeleteUserHandler)
	}
}
