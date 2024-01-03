package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ctrl "github.com/woonmapao/user-service-go/controllers"
	i "github.com/woonmapao/user-service-go/initializer"
	m "github.com/woonmapao/user-service-go/models"
	r "github.com/woonmapao/user-service-go/responses"
	v "github.com/woonmapao/user-service-go/validations"
)

// Handle the creation of a new user
func AddUserHandler(c *gin.Context) {

	// Get data from request body
	var body m.UserRequest
	err := ctrl.BindAndValidate(c, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := ctrl.StartTrx(c)
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

	err = ctrl.AddUser(&body, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = ctrl.CommitTrx(c, tx)
	if err != nil {
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.CreateSuccess(),
	)
}

// Retrieve a specific user based on their ID
func GetUserHandler(c *gin.Context) {

	id, err := ctrl.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	user, err := ctrl.GetUser(id, i.DB)
	if err != nil {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.GetSuccess(user),
	)
}

// Fetch a list of all users from the database
func GetUsersHandler(c *gin.Context) {

	users, err := ctrl.GetUsers(i.DB)
	if err != nil {
		c.JSON(http.StatusNotFound,
			r.GetError(c.Errors.Errors()))
		return
	}

	c.JSON(http.StatusOK,
		r.GetUsersSuccess(*users),
	)
}

// Handle the update of an existing user
func UpdateUserHandler(c *gin.Context) {

	id, err := ctrl.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Get data from request body
	var updData m.UserRequest
	err = ctrl.BindAndValidate(c, &updData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := ctrl.StartTrx(c)
	if err != nil {
		return
	}

	// Find the updating user (validation)
	_, err = ctrl.GetUser(id, tx)
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

	err = ctrl.UpdateUser(&updData, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = ctrl.CommitTrx(c, tx)
	if err != nil {
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.UpdateSuccess(),
	)
}

// DeleteUser deletes a user based on their ID
func DeleteUserHandler(c *gin.Context) {

	id, err := ctrl.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := ctrl.StartTrx(c)
	if err != nil {
		return
	}

	// Find the updating user (validation)
	_, err = ctrl.GetUser(id, tx)
	if err != nil {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	err = ctrl.DeleteUser(id, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = ctrl.CommitTrx(c, tx)
	if err != nil {
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		r.DeleteSuccess(),
	)
}
