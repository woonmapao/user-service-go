package responses

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/user-service-go/models"
)

func CreateSuccess() gin.H {
	return gin.H{
		"status":  "success",
		"message": "user added successfully",
	}
}

func DeleteSuccess() gin.H {
	return gin.H{
		"status":  "success",
		"message": "user deleted successfully",
	}
}

func GetSuccess(user *models.User) gin.H {
	return gin.H{
		"status":  "success",
		"message": "user fetched successfully",
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

func GetError(errors []string) gin.H {
	return gin.H{
		"status":  "error",
		"message": "failed to fetch",
		"data": gin.H{
			"errors": errors,
		},
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
