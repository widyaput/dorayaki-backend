package models

import (
	"fmt"
	"net/http"
	"reflect"
)

// TokoDorayaki represents Dorayaki shop in database.

type TokoDorayaki struct {
	TokoID     int64 `gorm:"primaryKey;ForeignKey:id;References:id;constraint:OnDelete:CASCADE" json:"toko_id"`
	DorayakiID int64 `gorm:"primaryKey;ForeignKey:id;References:id;constraint:OnDelete:CASCADE" json:"dorayaki_id"`
	Stok       int64 `gorm:"default:0;check:stok>=0" json:"stok"`
	CreatedAt  int64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  int64 `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns table's name inside database.
func (TokoDorayaki) TableName() string {
	return "toko_dorayaki"
}

type StokDorayaki struct {
	DorayakiID        int64  `json:"dorayaki_id"`
	DorayakiRasa      string `json:"dorayaki_rasa"`
	DorayakiDeskripsi string `json:"dorayaki_deskripsi"`
	DorayakiImageURL  string `json:"dorayaki_image_url"`
	Stok              int64  `json:"stok"`
}

type InputStok struct {
	AddStok int64 `json:"add_stok"`
}

type InputTransfer struct {
	Stock      int64 `json:"stock"`
	IdDorayaki int64 `json:"id_dorayaki"`
}

func (t *InputTransfer) Bind(r *http.Request) error {
	if t.Stock <= 0 {
		return fmt.Errorf("invalid stock, must greater than 0")
	}
	if t.IdDorayaki <= 0 {
		return fmt.Errorf("invalid id dorayaki")
	}
	return nil
}

// Reject addstock = 0
func (s *InputStok) Bind(r *http.Request) error {
	if reflect.TypeOf(s.AddStok).Kind().String() != "int" &&
		reflect.TypeOf(s.AddStok).Kind().String() != "int64" &&
		reflect.TypeOf(s.AddStok).Kind().String() != "int32" &&
		s.AddStok == 0 {
		return fmt.Errorf("invalid added stock")
	}
	return nil
}

// Bind accept body request and turn it into TokoDorayaki.
// Stok in body req means adding stok.
func (t *TokoDorayaki) Bind(r *http.Request) error {
	if t.TokoID <= 0 {
		return fmt.Errorf("toko id is required")
	}
	if t.DorayakiID <= 0 {
		return fmt.Errorf("dorayaki id is required")
	}
	return nil
}
