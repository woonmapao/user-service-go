package validations

import (
	"github.com/woonmapao/user-service-go/models"
	"gorm.io/gorm"
)

func IsUsernameDuplicate(username string, tx *gorm.DB) bool {

	var user models.User
	err := tx.Where("username = ?", username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// An unexpected error occurred, return true to handle it outside the function
		return true
	}

	// If a user with the given username is found, return true (duplicate)
	return user.ID != 0
}

func IsEmailDuplicate(email string, tx *gorm.DB) bool {

	var user models.User
	err := tx.Where("email = ?", email).First(&user).Error
	if err != nil {
		return false
	}

	// If a user with the given email is found, return true (duplicate)
	return user.ID != 0
}
