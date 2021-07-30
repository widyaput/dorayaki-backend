package seeds

import (
	"dorayaki/internal/helpers"
	"dorayaki/internal/models"

	"gorm.io/gorm"
)

func CreateCredentials(db *gorm.DB, username, password string) error {
	hashPwd, _ := helpers.HashPassword(password)
	return db.Create(&models.Credentials{Username: username, Password: hashPwd}).Error
}
