package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/user-service-go/controllers"
)

func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("/", controllers.GetAllUsers)
		userGroup.POST("/", controllers.AddUser)

		userGroup.GET("/:id", controllers.GetUserByID)
		userGroup.PUT("/:id", controllers.UpdateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)

		userGroup.GET("/:id/orders", controllers.GetUserOrders)
	}

}
