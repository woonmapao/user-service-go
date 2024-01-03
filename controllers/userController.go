package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/user-service-go/initializer"
	"github.com/woonmapao/user-service-go/models"
	r "github.com/woonmapao/user-service-go/responses"
	"github.com/woonmapao/user-service-go/validations"
	"gorm.io/gorm"
)

func AddUser(c *gin.Context) {
	// Handle the creation of a new user

	// Get data from the request body
	var body struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				"Invalid request format",
				err.Error(),
			}))
		return
	}

	// Check for empty values
	if body.Username == "" || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				"Username, email, and password are required fields",
			}))
		return
	}

	// Start a transaction
	tx, err := startTrx(c)
	if err != nil {
		return
	}

	// Check for duplicate username
	if validations.IsUsernameDuplicate(body.Username, tx) {
		tx.Rollback()
		c.JSON(http.StatusConflict,
			r.CreateError([]string{
				"Username is already taken",
			}))
		return
	}

	// Check for duplicate email
	if validations.IsEmailDuplicate(body.Email, tx) {
		tx.Rollback()
		c.JSON(http.StatusConflict,
			r.CreateError([]string{
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
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to create user",
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = commitTrx(c, tx)
	if err != nil {
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.CreateSuccessResponse(&user),
	)

}

// Retrieve a specific user based on their ID
func GetUserByID(c *gin.Context) {

	id, err := getID(c)
	if err != nil {
		return
	}

	// Get the user from the database
	var user models.User
	err = initializer.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to fetch user",
				err.Error(),
			}))
		return
	}
	// Check if the user was not found
	if user == (models.User{}) {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				"User not found",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.GetSuccessResponse(&user),
	)

}

// Fetch a list of all users from the database
func GetAllUsers(c *gin.Context) {

	// Get all users from the database
	var users []models.User
	err := initializer.DB.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to fetch users",
				err.Error(),
			}))
		return
	}
	// Check if no users were found
	if len(users) == 0 {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				"No users found",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.GetSuccessResponseForMultipleUsers(users),
	)

}

// Handle the update of an existing user
func UpdateUser(c *gin.Context) {

	id, err := getID(c)
	if err != nil {
		return
	}

	// Get data from request body
	var updateData struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				"Invalid request format",
				err.Error(),
			}))
		return
	}

	// Check for empty values
	if updateData.Username == "" || updateData.Email == "" || updateData.Password == "" {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				"Username, email, and password are required fields",
			}))
		return
	}

	// Start a transaction
	tx, err := startTrx(c)
	if err != nil {
		return
	}

	// Check if the user with the given ID exists
	var user models.User
	err = tx.First(&user, id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to fetch user",
				err.Error(),
			}))
		return
	}
	if user == (models.User{}) {
		tx.Rollback()
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				"User not found",
			}))
		return
	}

	// Check for duplicate username
	if validations.IsUsernameDuplicate(updateData.Username, tx) {
		tx.Rollback()
		c.JSON(http.StatusConflict,
			r.CreateError([]string{
				"Username is already taken",
			}))
		return
	}

	// Check for duplicate email
	if validations.IsEmailDuplicate(updateData.Email, tx) {
		tx.Rollback()
		c.JSON(http.StatusConflict,
			r.CreateError([]string{
				"Email is already registered",
			}))
		return
	}

	// Update user fields
	user.Username = updateData.Username
	user.Email = updateData.Email
	user.Password = updateData.Password

	// Save the updated user to the database
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to update user",
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	if err != nil {
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.UpdateSuccessResponse(&user),
	)

}

// GetUserOrders fetches all orders associated with a specific user
func GetUserOrders(c *gin.Context) {

	id, err := getID(c)
	if err != nil {
		return
	}

	// Get the user from the database
	var user models.User
	err = initializer.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to fetch user",
				err.Error(),
			}))
		return
	}
	// Check if the user was not found
	if user == (models.User{}) {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				"User not found",
			}))
		return
	}

	// Fetch orders for the user from order service

	url := "http://will-decide-later/api/orders?userId="
	url += strconv.Itoa(id)

	// Make HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to fetch user orders",
				err.Error(),
			}))
		return
	}

	if res.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to fetch user orders",
			}))
		return
	}

	// Decode the JSON response into OrderResponse struct
	var orderResponse models.OrderResponse
	err = json.NewDecoder(res.Body).Decode(&orderResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to fetch user orders",
				err.Error(),
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.CreateSuccessResponseForUserOrders(
			orderResponse.Data.Orders,
		),
	)

}

// DeleteUser deletes a user based on their ID
func DeleteUser(c *gin.Context) {

	id, err := getID(c)
	if err != nil {
		return
	}

	// Start a transaction
	tx, err := startTrx(c)
	if err != nil {
		return
	}

	// Check if the user with the given ID exists
	var user models.User
	err = tx.First(&user, id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to fetch user",
				err.Error(),
			}))
		return
	}
	if user == (models.User{}) {
		tx.Rollback()
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				"User not found",
			}))
		return
	}

	// Delete the user
	err = tx.Delete(&models.User{}, id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to delete user",
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	if err != nil {
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.DeleteSuccessResponse(&user),
	)
}

func getID(c *gin.Context) (id int, err error) {

	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				"Invalid user ID",
				err.Error(),
			}))
		return 0, err
	}

	return id, nil
}

func startTrx(c *gin.Context) (*gorm.DB, error) {

	tx := initializer.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to begin transaction",
				tx.Error.Error(),
			}))
		return nil, tx.Error
	}
	return tx, nil
}

func commitTrx(c *gin.Context, tx *gorm.DB) error {

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to commit transaction",
				err.Error(),
			}))
		return err
	}
	return nil
}

// Possible plan

// User Authentication (AuthenticateUser):

//     Authenticate users during login.
//     Verify provided credentials against stored user information.

// User Authorization (AuthorizeUser):

//     Determine whether a user has the necessary permissions to perform certain actions.

// User Search (SearchUsers):

//     Implement a search functionality for users based on criteria like name, email, etc.
