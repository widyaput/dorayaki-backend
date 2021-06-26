package models

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
