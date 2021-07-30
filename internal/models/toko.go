package models

import (
	"fmt"
	"net/http"
)

// Toko represents shop in database.

type Toko struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama      string         `gorm:"not null" json:"nama"`
	Jalan     string         `gorm:"not null" json:"jalan"`
	Kecamatan string         `gorm:"not null" json:"kecamatan"`
	Provinsi  string         `gorm:"not null" json:"provinsi"`
	Dorayaki  []Dorayaki     `gorm:"many2many:toko_dorayaki;constraint:OnDelete:CASCADE" json:"dorayaki"`
	Stok      []TokoDorayaki `json:"stok"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns table's name inside database.
func (Toko) TableName() string {
	return "toko"
}

// Bind accept body request and turn it into Toko
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
	// Optional
	for _, dora := range t.Dorayaki {
		if err := dora.Bind(r); err != nil {
			fmt.Println(err.Error())
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}
