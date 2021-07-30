package seeds

import (
	"dorayaki/internal/models"

	"gorm.io/gorm"
)

func CreateShop(db *gorm.DB, nama, jalan, kecamatan, provinsi string) error {
	return db.Create(&models.Toko{Nama: nama, Jalan: jalan, Kecamatan: kecamatan, Provinsi: provinsi}).Error
}
