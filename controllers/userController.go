package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	i "github.com/woonmapao/user-service-go/initializer"
	"github.com/woonmapao/user-service-go/models"
	r "github.com/woonmapao/user-service-go/responses"
	v "github.com/woonmapao/user-service-go/validations"
	"gorm.io/gorm"
)

// Handle the creation of a new user
func AddUser(c *gin.Context) {

	// Get data from request body
	var body models.UserRequest
	err := bindAndValidate(c, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := startTrx(c)
	if err != nil {
		return
	}

	// Check for duplicate username and email
	dupe, err := v.IsDupe(body.Username, body.Email, tx)
	if err != nil { // Failed to fetch
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}
	if dupe { // Found duplicate
		c.JSON(http.StatusConflict,
			r.CreateError([]string{
				"found duplicate record",
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
				"failed to create user",
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
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Get the user from the database
	var user models.User
	err = i.DB.First(&user, id).Error
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

	var users []models.User
	err := i.DB.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"failed to fetch users",
			}))
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				"no users found",
			}))
		return
	}

	c.JSON(http.StatusOK,
		r.GetUsersSuccess(users),
	)

}

// Handle the update of an existing user
func UpdateUser(c *gin.Context) {

	id, err := getID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Get data from request body
	var updData models.UserRequest
	err = bindAndValidate(c, &updData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := startTrx(c)
	if err != nil {
		return
	}

	// Find the updating user
	user, err := getUser(id, tx)
	if err != nil {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Check for duplicate username and email
	dupe, err := v.IsDupe(updData.Username, updData.Email, tx)
	if err != nil { // Failed to fetch
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}
	if dupe { // Found duplicate
		c.JSON(http.StatusConflict,
			r.CreateError([]string{
				"found duplicate record",
			}))
		return
	}

	// Update user fields
	user.Username = updData.Username
	user.Email = updData.Email
	user.Password = updData.Password

	// Save the updated user to the database
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"failed to update user",
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
		r.UpdateSuccess(),
	)

}

// GetUserOrders fetches all orders associated with a specific user
func GetUserOrders(c *gin.Context) {

	id, err := getID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Get the user from the database
	var user models.User
	err = i.DB.First(&user, id).Error
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
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
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

func getID(c *gin.Context) (int, error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errors.New("invalid user id")
	}
	return id, nil
}

func startTrx(c *gin.Context) (*gorm.DB, error) {

	tx := i.DB.Begin()
	if tx.Error != nil {
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

func bindAndValidate(c *gin.Context, body *models.UserRequest) error {

	err := c.ShouldBindJSON(&body)
	if err != nil {
		return errors.New(
			"invalid request format",
		)
	}
	if body.Username == "" ||
		body.Email == "" ||
		body.Password == "" {
		return errors.New(
			"username, email and password are required fields",
		)
	}
	return nil
}

func getUser(id int, tx *gorm.DB) (*models.User, error) {

	var user models.User
	err := tx.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return &user, errors.New("user not found")
	}
	if err != nil {
		return &user, errors.New("something went wrong")
	}
	return &user, nil
}

// Possible plan

// User Authentication (AuthenticateUser):

//     Authenticate users during login.
//     Verify provided credentials against stored user information.

// User Authorization (AuthorizeUser):

//     Determine whether a user has the necessary permissions to perform certain actions.

// User Search (SearchUsers):

//     Implement a search functionality for users based on criteria like name, email, etc.
