package models

import (
	"fmt"
	"net/http"
	"net/url"
)

// Dorayaki represents dorayaki in database.
type Dorayaki struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	Rasa      string `gorm:"not null"`
	Deskripsi string `gorm:"not null"`
	ImageURL  string `gorm:"not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

// TableName returns table's name inside database.
func (Dorayaki) TableName() string {
	return "dorayaki"
}

func (d *Dorayaki) Bind(r *http.Request) error {
	if d.Rasa == "" {
		return fmt.Errorf("rasa is required")
	}
	if d.Deskripsi == "" {
		return fmt.Errorf("deskripsi is required")
	}
	// if d.Base64 == "" {
	// 	return fmt.Errorf("base64 image is required")
	// }
	if _, err := url.ParseRequestURI(d.ImageURL); err != nil {
		return fmt.Errorf("imageurl should be valid")
	}
	return nil
}

func (*Dorayaki) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
