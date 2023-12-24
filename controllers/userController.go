package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/user-service-go/initializer"
	"github.com/woonmapao/user-service-go/models"
	"github.com/woonmapao/user-service-go/responses"
	"github.com/woonmapao/user-service-go/validations"
)

func AddUser(c *gin.Context) {
	// Handle the creation of a new user

	// Get data from the request body
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validations.IsUsernameDuplicate(body.Username) {
		c.JSON(http.StatusConflict,
			responses.CreateErrorResponse([]string{
				"Username is already taken",
			}))
		return
	}

	// Check for duplicate email
	if validations.IsEmailDuplicate(body.Email) {
		c.JSON(http.StatusConflict,
			responses.CreateErrorResponse([]string{
				"Email is already registered",
			}))
		return
	}

	// Create user in the database
	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	}

	err = initializer.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to create user",
			}))
		return
	}

	// Return status
	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&user),
	)

}

func GetUserByID(c *gin.Context) {
	// Retrieve a specific user based on their ID

	// Get ID from URL param
	id := c.Param("id")

	// Get the user from the database
	var user models.User
	err := initializer.DB.First(&user, id).Error
	if err != nil {
		// Handle user not found or other errors
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Respond with the found user
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}

// Fetch a list of all users from the database
func GetAllUsers(c *gin.Context) {

	// Get all users from the database
	var users []models.User
	err := initializer.DB.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	// Respond with the found users
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})

}

func UpdateUser(c *gin.Context) {
	// Handle the update of an existing user

	// Get ID from URL param
	id := c.Param("id")

	// Get data from request body
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the user with the given ID exists
	var user models.User
	err = initializer.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Update user
	initializer.DB.Model(&user).Updates(models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	})

	// Respond with the updated user
	c.JSON(http.StatusOK, gin.H{
		"updatedUser": user,
	})

}

// GetUserOrders fetches all orders associated with a specific user
func GetUserOrders(c *gin.Context) {
	// Extract user ID from the request parameters
	userID := c.Param("id")

	// Query the database for orders associated with the user
	var userOrders []models.Order
	if err := initializer.DB.Where("user_id = ?", userID).Find(&userOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user orders",
		})
		return
	}

	// Return a JSON response with the user's orders
	c.JSON(http.StatusOK, gin.H{
		"user_orders": userOrders,
	})
}

// DeleteUser deletes a user based on their ID
func DeleteUser(c *gin.Context) {
	// Get the ID off the URL
	id := c.Param("id")

	// Delete the user
	err := initializer.DB.Delete(&models.User{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
