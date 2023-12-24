package validations

import (
	"github.com/woonmapao/user-service-go/initializer"
	"github.com/woonmapao/user-service-go/models"
)

func IsUsernameDuplicate(username string) bool {
	var count int64
	initializer.DB.Model(&models.User{}).Where(
		"username = ?", username).Count(&count)
	return count > 0
}

func IsEmailDuplicate(email string) bool {
	var count int64
	initializer.DB.Model(&models.User{}).Where(
		"email = ?", email).Count(&count)
	return count > 0
}
