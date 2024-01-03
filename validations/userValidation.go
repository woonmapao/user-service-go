package validations

import (
	"errors"

	"github.com/woonmapao/user-service-go/models"
	"gorm.io/gorm"
)

// Check username and email dupe
func IsDupe(username, email string, tx *gorm.DB) (bool, error) {

	// Check username dupe
	var a models.User
	err := tx.Where("username = ?", username).First(&a).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, errors.New("failed to validate")
	}
	// Check email dupe
	var b models.User
	err = tx.Where("email = ?", email).First(&b).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, errors.New("failed to validate")
	}
	return a.ID != 0 || b.ID != 0, nil // found dupe
}
