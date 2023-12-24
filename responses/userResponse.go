package responses

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/user-service-go/models"
)

func CreateSuccessResponse(user *models.User) gin.H {
	return gin.H{
		"status":  "success",
		"message": "User added successfully",
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	}
}

func CreateErrorResponse(errors []string) gin.H {
	return gin.H{
		"status":  "error",
		"message": "Validation failed",
		"data": gin.H{
			"errors": errors,
		},
	}
}
