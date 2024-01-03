package validations

import (
	"errors"

	"github.com/woonmapao/user-service-go/models"
	"gorm.io/gorm"
)

// Check username and email dupe
func IsDupe(username, email string, tx *gorm.DB) (bool, error) {

	// Check username dupe
	var user models.User
	err := tx.Where("username = ?", username).First(&user).Error
	if err != nil {
		// failed to fetch
		return false, errors.New("failed to validate")
	}
	if user != (models.User{}) {
		return true, nil // found dupe
	}

	// Check email dupe
	err = tx.Where("email = ?", email).First(&user).Error
	if err != nil {
		// failed to fetch
		return false, errors.New("failed to validate")
	}
	if user != (models.User{}) {
		return true, nil // found dupe
	}
	return false, nil // no username or email dupe
}
