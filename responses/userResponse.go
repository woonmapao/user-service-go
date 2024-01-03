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

func DeleteSuccessResponse(user *models.User) gin.H {
	return gin.H{
		"status":  "success",
		"message": "User deleted successfully",
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	}
}

func GetSuccessResponse(user *models.User) gin.H {
	return gin.H{
		"status":  "success",
		"message": "User fetched successfully",
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	}
}

func UpdateSuccess() gin.H {
	return gin.H{
		"status":  "success",
		"message": "user updated successfully",
	}
}

func CreateError(errors []string) gin.H {
	return gin.H{
		"status":  "error",
		"message": "Validation failed",
		"data": gin.H{
			"errors": errors,
		},
	}
}

func GetUsersSuccess(userList []models.User) gin.H {

	users := make([]gin.H, len(userList))
	for i, user := range userList {
		users[i] = map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		}
	}
	return gin.H{
		"status":  "success",
		"message": "users fetched successfully",
		"data": gin.H{
			"users": users,
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
