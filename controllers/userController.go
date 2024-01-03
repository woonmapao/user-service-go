package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	i "github.com/woonmapao/user-service-go/initializer"
	m "github.com/woonmapao/user-service-go/models"
	r "github.com/woonmapao/user-service-go/responses"
	"gorm.io/gorm"
)

func AddUser(user *m.UserRequest, tx *gorm.DB) error {

	adding := m.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	err := tx.Create(&adding).Error
	if err != nil {
		tx.Rollback()
		return errors.New("failed to create user")
	}
	return nil
}

func GetUser(id int, db *gorm.DB) (*m.User, error) {

	var user m.User
	err := db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return &user, errors.New("user not found")
	}
	if err != nil {
		return &user, errors.New("something went wrong")
	}
	return &user, nil
}

func GetUsers(db *gorm.DB) (*[]m.User, error) {

	var users []m.User
	err := db.Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &users, errors.New("failed to fetch user")
	}
	if err == gorm.ErrRecordNotFound {
		return &users, errors.New("no user found")
	}
	return &users, nil
}

func UpdateUser(user *m.UserRequest, tx *gorm.DB) error {

	updating := m.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	err := tx.Save(&updating).Error
	if err != nil {
		tx.Rollback()
		return errors.New("failed to update user")
	}
	return nil
}

func DeleteUser(id int, tx *gorm.DB) error {

	err := tx.Delete(&m.User{}, id).Error
	if err != nil {
		tx.Rollback()
		return errors.New("failed to delete user")
	}
	return nil
}

func GetID(c *gin.Context) (int, error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errors.New("invalid user id")
	}
	return id, nil
}

func StartTrx(c *gin.Context) (*gorm.DB, error) {

	tx := i.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func CommitTrx(c *gin.Context, tx *gorm.DB) error {

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

func BindAndValidate(c *gin.Context, body *m.UserRequest) error {

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
