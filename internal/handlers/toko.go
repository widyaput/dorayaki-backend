package handlers

import (
	"context"
	"dorayaki/configs/database"
	"dorayaki/internal/helpers"
	"dorayaki/internal/models"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	keyShop       key = iota
	keyShopTarget key = iota
	keyDorayaki   key = iota
)

func shops(router chi.Router) {
	router.Group(func(router chi.Router) {
		router.Get("/", getAllShop)
		router.Get("/search", paginateShop)
		router.Group(func(router chi.Router) {
			router.Use(authenticator)
			router.Post("/", createShop)
		})
		router.Route("/{shopId}", func(router chi.Router) {
			router.Use(ShopContext)
			router.Get("/", getShop)
			router.Group(func(router chi.Router) {
				router.Use(authenticator)
				router.Put("/", updateShop)
				router.Delete("/", deleteShop)
			})
			router.Route("/{dorayakiId}", func(router chi.Router) {
				router.Use(DorayakiContext)
				router.Get("/", getStok)
				router.Group(func(router chi.Router) {
					router.Use(authenticator)
					router.Post("/", addStok)
				})
			})
			router.Route("/transfer/{targetShopId}", func(router chi.Router) {
				router.Use(TargetShopContext)
				router.Use(authenticator)
				router.Post("/", transferStok)
			})
		})
	})
}

// createShop will create shop, omit all dorayakis.
func createShop(w http.ResponseWriter, r *http.Request) {
	var shop models.Toko
	if err := render.Bind(r, &shop); err != nil {
		render.Render(w, r, models.ErrBadRequest)
		return
	}
	// Omit dorayaki. Dorayaki need to added to database using /api/v1/dorayakis.
	if rs := database.DB.Omit("Dorayaki").Omit("Stok").Create(&shop); rs.Error != nil {
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

// updateShop will update information about specific shop. Will not update it's dorayaki.
func updateShop(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(keyShop).(int)
	var newShop models.Toko
	var oldShop models.Toko
	if err := render.Bind(r, &newShop); err != nil {
		render.Render(w, r, models.ErrBadRequest)
		return
	}
	if rs := database.DB.Where("id = ?", id).Preload("Dorayaki").Preload("Stok").
		First(&oldShop); rs.Error != nil {
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

// deleteShop will delete shop with specific id.
func deleteShop(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(keyShop).(int)
	var oldShop models.Toko
	if rs := database.DB.Where("ID = ?", id).First(&oldShop); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	if rs := database.DB.Delete(&models.Toko{}, id); rs.Error != nil {
		render.Render(w, r, models.ServerErrorRenderer(rs.Error))
		return
	}
}

// getShop will retrive specific shop by id and retrieve all dorayaki it has.
func getShop(w http.ResponseWriter, r *http.Request) {
	var shop models.Toko
	id := r.Context().Value(keyShop).(int)
	if rs := database.DB.Where("id = ?", id).Preload("Dorayaki").Preload("Stok").
		First(&shop); rs.Error != nil {
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

// getAllShop will get all shop without retrieving dorayaki it has.
func getAllShop(w http.ResponseWriter, r *http.Request) {
	var list []models.Toko
	if rs := database.DB.Find(&list); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	resp := models.ResponseToko{Response: *models.SuccessResponse}
	resp.Data = list
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

// getStok will get stock of specific dorayaki at specific shop.
func getStok(w http.ResponseWriter, r *http.Request) {
	idShop := r.Context().Value(keyShop).(int)
	idDorayaki := r.Context().Value(keyDorayaki).(int)
	var stok models.TokoDorayaki
	if rs := database.DB.
		FirstOrCreate(&stok,
			models.TokoDorayaki{
				TokoID:     int64(idShop),
				DorayakiID: int64(idDorayaki),
			}); rs.Error != nil {

		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	resp := models.ResponseStok{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, stok)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

// addStok will add specific dorayaki at specific shop.
// Can reduce stok with input stock has negative value.
func addStok(w http.ResponseWriter, r *http.Request) {
	idShop := r.Context().Value(keyShop).(int)
	idDorayaki := r.Context().Value(keyDorayaki).(int)
	var stok models.TokoDorayaki
	var addStock models.InputStok
	if err := render.Bind(r, &addStock); err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
	if rs := database.DB.
		FirstOrCreate(&stok,
			models.TokoDorayaki{
				TokoID:     int64(idShop),
				DorayakiID: int64(idDorayaki),
			}); rs.Error != nil {

		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	if stok.Stok+addStock.AddStok < 0 {
		render.Render(w, r, models.ErrorRenderer(fmt.Errorf("stok tidak mencukupi")))
		return
	}
	stok.Stok += addStock.AddStok
	if rs := database.DB.Save(&stok); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	resp := models.ResponseStok{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, stok)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

// transferStok will substitute some dorayakis from source shop and transfer it to target shop.
// input stock must be positive.
func transferStok(w http.ResponseWriter, r *http.Request) {
	idSource := r.Context().Value(keyShop).(int)
	idTarget := r.Context().Value(keyShopTarget).(int)
	var source models.TokoDorayaki
	var target models.TokoDorayaki
	var stock models.InputTransfer
	if err := render.Bind(r, &stock); err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
	if rs := database.DB.
		First(&source,
			models.TokoDorayaki{
				TokoID:     int64(idSource),
				DorayakiID: int64(stock.IdDorayaki),
			}); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	if source.Stok-stock.Stock < 0 {
		render.Render(w, r, models.ErrorRenderer(fmt.Errorf("stock tidak mencukupi")))
		return
	}
	if rs := database.DB.
		FirstOrCreate(&target,
			models.TokoDorayaki{
				TokoID:     int64(idTarget),
				DorayakiID: int64(stock.IdDorayaki),
			}); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	target.Stok += stock.Stock
	source.Stok -= stock.Stock
	if rs := database.DB.Save(&target); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	if rs := database.DB.Save(&source); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	if err := render.Render(w, r, models.SuccessResponse); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func paginateShop(w http.ResponseWriter, r *http.Request) {
	rawQuery, rawArgs, err := helpers.PaginateAbstract(models.Toko{}.TableName(), r)
	if err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
	totalRawQuery, totalRawArgs, err := helpers.TakeQuery(models.Toko{}.TableName(), r).ToSql()
	if err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
	var data []models.Toko
	if rs := database.DB.Raw(rawQuery, rawArgs...).Scan(&data); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	if data == nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	var totalData []models.Toko
	if rs := database.DB.Raw(totalRawQuery, totalRawArgs...).Scan(&totalData); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	sort := r.URL.Query().Get("sort")
	var idxPage int
	var itemsPerPage int
	idxPage, err = strconv.Atoi(r.URL.Query().Get("pageIndex"))
	if err != nil {
		idxPage = 1
	}
	itemsPerPage, err = strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	if err != nil {
		itemsPerPage = 10
	}
	respPaginate := models.ResponsePaginate{
		Response:     *models.SuccessResponse,
		ItemsPerPage: int64(itemsPerPage),
		TotalItems:   int64(len(totalData)),
		PageIndex:    int64(idxPage),
		TotalPages:   int64(math.Ceil(float64(len(totalData)) / float64(itemsPerPage))),
		Sort:         sort,
	}
	resp := models.ResponsePaginateToko{
		ResponsePaginate: respPaginate,
		Kecamatan:        r.URL.Query().Get("kecamatan"),
		Provinsi:         r.URL.Query().Get("provinsi"),
		Data:             data,
	}
	if err = render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
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
func TargetShopContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shopId := chi.URLParam(r, "targetShopId")
		if shopId == "" {
			render.Render(w, r, models.ErrorRenderer(fmt.Errorf("target shop ID is required")))
			return
		}
		id, err := strconv.Atoi(shopId)
		if err != nil {
			render.Render(w, r, models.ErrorRenderer(fmt.Errorf("invalid target shop ID")))
			return
		}
		ctx := context.WithValue(r.Context(), keyShopTarget, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
