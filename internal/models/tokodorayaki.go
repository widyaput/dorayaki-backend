package models

import (
	"fmt"
	"net/http"
	"reflect"
)

// TokoDorayaki represents Dorayaki shop in database.

type TokoDorayaki struct {
	TokoID     int64 `gorm:"primaryKey;ForeignKey:id;References:id;constraint:OnDelete:CASCADE"`
	DorayakiID int64 `gorm:"primaryKey;ForeignKey:id;References:id;constraint:OnDelete:CASCADE"`
	Stok       int64 `gorm:"default:0;check:stok>=0"`
	CreatedAt  int64 `gorm:"autoCreateTime"`
	UpdatedAt  int64 `gorm:"autoUpdateTime"`
}

// TableName returns table's name inside database.
func (TokoDorayaki) TableName() string {
	return "toko_dorayaki"
}

type InputStok struct {
	AddStok int64 `json:"add_stok"`
}

func (s *InputStok) Bind(r *http.Request) error {
	if reflect.TypeOf(s.AddStok).Kind().String() != "int" &&
		reflect.TypeOf(s.AddStok).Kind().String() != "int64" && reflect.TypeOf(s.AddStok).Kind().String() != "int32" {
		return fmt.Errorf("invalid added stock")
	}
	return nil
}

// Bind accept body request and turn it into TokoDorayaki.
// Stok in body req means adding stok.
func (t *TokoDorayaki) Bind(r *http.Request) error {
	if t.TokoID <= 0 {
		return fmt.Errorf("Toko id is required")
	}
	if t.DorayakiID <= 0 {
		return fmt.Errorf("Dorayaki id is required")
	}
	return nil
}
