package models

import (
	"fmt"
	"net/http"
)

// Toko represents shop in database.
type Toko struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	Nama      string     `gorm:"not null"`
	Jalan     string     `gorm:"not null"`
	Kecamatan string     `gorm:"not null"`
	Provinsi  string     `gorm:"not null"`
	Dorayaki  []Dorayaki `gorm:"many2many:toko_dorayaki"`
	CreatedAt int64      `gorm:"autoCreateTime"`
	UpdatedAt int64      `gorm:"autoUpdateTime"`
}

// TableName returns table's name inside database.
func (Toko) TableName() string {
	return "toko"
}

func (t *Toko) Bind(r *http.Request) error {
	if t.Nama == "" {
		return fmt.Errorf("nama is required")
	}
	if t.Jalan == "" {
		return fmt.Errorf("jalan is required")
	}
	if t.Kecamatan == "" {
		return fmt.Errorf("kecamatan is required")
	}
	if t.Provinsi == "" {
		return fmt.Errorf("provinsi is required")
	}
	return nil
}
