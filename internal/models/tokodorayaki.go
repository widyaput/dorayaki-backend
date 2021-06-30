package models

import "net/http"

// TokoDorayaki represents Dorayaki shop in database.
type TokoDorayaki struct {
	Toko       Toko
	Dorayaki   Dorayaki
	TokoID     int64 `gorm:"primaryKey"`
	DorayakiID int64 `gorm:"primaryKey"`
	Stok       int64 `gorm:"default:0;check:stok>=0"`
	CreatedAt  int64 `gorm:"autoCreateTime"`
	UpdatedAt  int64 `gorm:"autoUpdateTime"`
}

// TableName returns table's name inside database.
func (TokoDorayaki) TableName() string {
	return "toko_dorayaki"
}

func (t *TokoDorayaki) Bind(r *http.Request) error {
	return nil
}
