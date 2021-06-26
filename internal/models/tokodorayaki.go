package models

// TokoDorayaki represents Dorayaki shop in database.
type TokoDorayaki struct {
	Toko       Toko
	Dorayaki   Dorayaki
	TokoID     int64 `gorm:"primaryKey"`
	DorayakiID int64 `gorm:"primaryKey"`
	Stok       int64 `gorm:"default:0;check:stok>=0"`
}

// TableName returns table's name inside database.
func (TokoDorayaki) TableName() string {
	return "toko_dorayaki"
}
