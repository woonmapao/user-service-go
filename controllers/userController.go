package controllers

import (
	"net/http"
	"strconv"

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

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&user),
	)

}

// Retrieve a specific user based on their ID
func GetUserByID(c *gin.Context) {

	// Get ID from URL param
	userID := c.Param("id")

	// Convert user ID to integer (validations)
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid user ID",
			}))
		return
	}
	// Get the user from the database
	var user models.User
	err = initializer.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch user",
			}))
		return
	}
	// Check if the user was not found
	if &user == nil {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"User not found",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&user))

}

// Fetch a list of all users from the database
func GetAllUsers(c *gin.Context) {

	// Get all users from the database
	var users []models.User
	err := initializer.DB.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch users",
			}))
		return
	}
	// Check if no users were found
	if len(users) == 0 {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"No users found",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponseForMultipleUsers(users))

}

func UpdateUser(c *gin.Context) {
	// Handle the update of an existing user

	// Get ID from URL param
	userID := c.Param("id")

	// Convert user ID to integer (validations)
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid user ID",
			}))
		return
	}

	// Get data from request body
	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
			}))
		return
	}

	// Check if the user with the given ID exists
	var user models.User
	err = initializer.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch user",
			}))
		return
	}
	if &user == nil {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"User not found",
			}))
		return
	}

	// Update user fields
	user.Username = updateData.Username
	user.Email = updateData.Email
	user.Password = updateData.Password

	// Save the updated user to the database
	err = initializer.DB.Save(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to update user",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&user))

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
	userID := c.Param("id")

	// Convert user ID to integer (validations)
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid user ID",
			}))
		return
	}

	// Check if the user with the given ID exists
	var user models.User
	err = initializer.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch user",
			}))
		return
	}
	if &user == nil {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"User not found",
			}))
		return
	}

	// Delete the user
	err = initializer.DB.Delete(&models.User{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to delete user",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(nil))
}
