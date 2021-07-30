package seeds

import (
	"dorayaki/internal/seed"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

var rasaDorayaki = []string{"Pempek", "Cokelat", "Pisang", "nasgor", "Stroberi", "Duren"}
var jalan = []string{"Jalan Kebenaran", "Jalan Sesama", "Jalanin aja dulu", "Jalan-jalan yuk", "Jalan yang sempit"}
var kecamatan = []string{"Wakanda", "Mranggen", "Tayu", "Klipang", "Dukuhseti", "Area 51", "Bedagan"}
var provinsi = []string{"Sunda Empire", "Mondstadt", "Inazuma", "Tokyo", "Paris", "Ohio", "Manila"}
var toko = []string{"Takoyaki", "Tiara", "Berlian", "Permata", "Emas", "Perak"}

const defaultIMG = "https://upload.wikimedia.org/wikipedia/commons/7/7f/Dorayaki_001.jpg"
const lorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut quam mauris, bibendum sed nibh sed, semper tincidunt lorem. Maecenas vestibulum, nisi egestas semper imperdiet, ipsum felis maximus velit, tempor pretium mauris ex a metus. Curabitur eu ante eu nisi faucibus volutpat. Curabitur ac semper dolor, sit amet auctor leo. Sed non sagittis lorem. Quisque vehicula euismod dapibus. Fusce in quam nec nisl faucibus aliquam at ut leo. Vestibulum non sem eu dui sollicitudin blandit eget at massa. Proin hendrerit odio non magna condimentum, congue consequat velit pharetra."

func All() []seed.Seed {
	seeds := []seed.Seed{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 101; i++ {
		seeds = append(seeds, seed.Seed{
			Name: "CreateShop " + fmt.Sprint(i),
			Run: func(db *gorm.DB) error {
				return CreateShop(
					db,
					toko[rand.Intn(len(toko))],
					jalan[rand.Intn(len(jalan))],
					kecamatan[rand.Intn(len(kecamatan))],
					provinsi[rand.Intn(len(provinsi))],
				)
			},
		})
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 15; i++ {
		desc := "Dorayaki ke " + fmt.Sprint(i) + `
		` + lorem
		seeds = append(seeds, seed.Seed{
			Name: "CreateDorayaki " + fmt.Sprint(i),
			Run: func(db *gorm.DB) error {
				return CreateDorayaki(
					db,
					rasaDorayaki[rand.Intn(len(rasaDorayaki))],
					desc,
					defaultIMG,
				)
			},
		})
	}
	seeds = append(seeds, seed.Seed{
		Name: "CreateAdmin",
		Run: func(db *gorm.DB) error {
			return CreateCredentials(db, "admin", "dorayakidev")
		},
	})
	return seeds
}
