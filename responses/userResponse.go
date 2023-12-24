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

func CreateSuccessResponseForMultipleUsers(users []models.User) gin.H {
	userList := make([]map[string]interface{}, len(users))

	for i, user := range users {
		userList[i] = map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		}
	}

	return gin.H{
		"status":  "success",
		"message": "Users fetched successfully",
		"data": gin.H{
			"users": userList,
		},
	}
}

// createSuccessResponseForUserOrders formats the success response for user orders
func CreateSuccessResponseForUserOrders(orders []models.Order) gin.H {
	return gin.H{
		"status":  "success",
		"message": "Orders retrieved successfully",
		"data": gin.H{
			"orders": orders,
		},
	}
}
