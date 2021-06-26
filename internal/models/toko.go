package models

// Toko represents shop in database.
type Toko struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	Nama      string `gorm:"not null"`
	Jalan     string `gorm:"not null"`
	Kecamatan string `gorm:"not null"`
	Provinsi  string `gorm:"not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

// TableName returns table's name inside database.
func (Toko) TableName() string {
	return "toko"
}
