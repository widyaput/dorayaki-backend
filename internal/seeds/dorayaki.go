package seeds

import (
	"dorayaki/internal/models"

	"gorm.io/gorm"
)

func CreateDorayaki(db *gorm.DB, rasa, deskripsi, image string) error {
	return db.Create(&models.Dorayaki{Rasa: rasa, Deskripsi: deskripsi, ImageURL: image}).Error
}
