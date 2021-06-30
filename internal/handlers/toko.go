package handlers

import (
	"context"
	"dorayaki/configs/database"
	"dorayaki/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	keyShop key = iota
)

func shops(router chi.Router) {
	router.Get("/", getAllShop)
	router.Post("/", createShop)
	// TODO: Pagination
	router.Route("/{shopId}", func(router chi.Router) {
		router.Use(ShopContext)
		router.Get("/", getShop)
		router.Put("/", updateShop)
		router.Delete("/", deleteShop)
	})
}

func createShop(w http.ResponseWriter, r *http.Request) {
	var shop models.Toko
	if err := render.Bind(r, &shop); err != nil {
		render.Render(w, r, models.ErrBadRequest)
		return
	}
	if rs := database.DB.Create(&shop); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	resp := models.ResponseToko{Response: *models.SuccessCreateResponse}
	resp.Data = append(resp.Data, shop)

	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func updateShop(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(keyShop).(int)
	var newShop models.Toko
	var oldShop models.Toko
	if err := render.Bind(r, &newShop); err != nil {
		render.Render(w, r, models.ErrBadRequest)
		return
	}
	if rs := database.DB.Where("ID = ?", id).First(&oldShop); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	oldShop.Jalan = newShop.Jalan
	oldShop.Nama = newShop.Nama
	oldShop.Kecamatan = newShop.Kecamatan
	oldShop.Provinsi = newShop.Provinsi
	if rs := database.DB.Save(&oldShop); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	resp := models.ResponseToko{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, oldShop)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func deleteShop(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(keyShop).(int)
	var oldShop models.Toko
	if rs := database.DB.Where("ID = ?", id).First(&oldShop); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	if rs := database.DB.Delete(&models.Toko{}, id); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	// resp := *models.SuccessDeleteResponse
	// if err := render.Render(w, r, &resp); err != nil {
	// 	render.Render(w, r, models.ServerErrorRenderer(err))
	// 	return
	// }
}

func getShop(w http.ResponseWriter, r *http.Request) {
	var shop models.Toko
	id := r.Context().Value(keyShop).(int)
	if rs := database.DB.Where("ID = ?", id).First(&shop); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	resp := models.ResponseToko{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, shop)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func getAllShop(w http.ResponseWriter, r *http.Request) {
	var list []models.Toko
	if rs := database.DB.Find(&list); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	resp := models.ResponseToko{Response: *models.SuccessResponse}
	resp.Data = list
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
}

func ShopContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shopId := chi.URLParam(r, "shopId")
		if shopId == "" {
			render.Render(w, r, models.ErrorRenderer(fmt.Errorf("shop ID is required")))
			return
		}
		id, err := strconv.Atoi(shopId)
		if err != nil {
			render.Render(w, r, models.ErrorRenderer(fmt.Errorf("invalid shop ID")))
			return
		}
		ctx := context.WithValue(r.Context(), keyShop, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
