package handlers

import (
	"context"
	"dorayaki/configs/database"
	"dorayaki/internal/models"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type key int64

const maxUploadSize = 2 * 1024 * 1024
const uploadPath = "internal/assets/"

func dorayakis(router chi.Router) {
	router.Group(func(router chi.Router) {
		router.Get("/", getAllDorayaki)
		router.Get("/search", paginateDorayakiGorm)
		router.Group(func(router chi.Router) {
			router.Use(authenticator)
			router.Post("/", createDorayaki)
		})
		router.Route("/{dorayakiId}", func(router chi.Router) {
			router.Use(DorayakiContext)
			router.Get("/", getDorayaki)
			router.Group(func(router chi.Router) {
				router.Use(authenticator)
				router.Put("/", updateDorayaki)
				router.Delete("/", deleteDorayaki)
				router.Post("/upload", uploadImageDorayaki)
			})
		})
	})
}

// createDorayaki based on body request
func createDorayaki(w http.ResponseWriter, r *http.Request) {
	var dorayaki models.Dorayaki
	if err := render.Bind(r, &dorayaki); err != nil {
		render.Render(w, r, models.ErrBadRequest)
		return
	}
	if rs := database.DB.Create(&dorayaki); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer((rs.Error)))
		return
	}
	resp := models.ResponseDorayaki{Response: *models.SuccessCreateResponse}
	resp.Data = append(resp.Data, dorayaki)

	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer((err)))
		return
	}
}

// updateDorayaki by id based on body request
func updateDorayaki(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(keyDorayaki).(int)
	var newDorayaki models.Dorayaki
	var oldDorayaki models.Dorayaki
	if err := render.Bind(r, &newDorayaki); err != nil {
		log.Print(err.Error())
		render.Render(w, r, models.ErrBadRequest)
		return
	}
	if rs := database.DB.Where("ID = ?", id).First(&oldDorayaki); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	oldDorayaki.Deskripsi = newDorayaki.Deskripsi
	oldDorayaki.Rasa = newDorayaki.Rasa
	if oldDorayaki.ImageURL != newDorayaki.ImageURL {
		nameOfImage := strings.ReplaceAll(oldDorayaki.ImageURL, "http://localhost:8080/api/v1/files/", "")
		if _, err := os.Stat(uploadPath + nameOfImage); err == nil {
			os.Remove(uploadPath + nameOfImage)
		}
	}
	oldDorayaki.ImageURL = newDorayaki.ImageURL
	if rs := database.DB.Save(&oldDorayaki); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer((rs.Error)))
		return
	}
	resp := models.ResponseDorayaki{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, oldDorayaki)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

// deleteDorayaki by id
func deleteDorayaki(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(keyDorayaki).(int)
	var oldDorayaki models.Dorayaki
	if rs := database.DB.Where("ID = ?", id).First(&oldDorayaki); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	nameOfImage := strings.ReplaceAll(oldDorayaki.ImageURL, "http://localhost:8080/api/v1/files/", "")
	log.Print(nameOfImage)
	if _, err := os.Stat(uploadPath + nameOfImage); err == nil {
		os.Remove(uploadPath + nameOfImage)
	}
	if rs := database.DB.Delete(&models.Dorayaki{}, id); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
}

// getDorayaki by id
func getDorayaki(w http.ResponseWriter, r *http.Request) {
	var dorayaki models.Dorayaki
	id := r.Context().Value(keyDorayaki).(int)
	if rs := database.DB.Where("ID = ?", id).First(&dorayaki); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	resp := models.ResponseDorayaki{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, dorayaki)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

// get all dorayaki
func getAllDorayaki(w http.ResponseWriter, r *http.Request) {
	var list []models.Dorayaki
	if rs := database.DB.Find(&list); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	resp := models.ResponseDorayaki{Response: *models.SuccessResponse}
	resp.Data = list
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
}

// uploadImage of dorayaki into filesystem
func uploadImageDorayaki(w http.ResponseWriter, r *http.Request) {
	var dorayaki models.Dorayaki
	id := r.Context().Value(keyDorayaki).(int)
	if rs := database.DB.Where("ID = ?", id).First(&dorayaki); rs.Error != nil {
		render.Render(w, r, models.ErrNotFound)
		return
	}
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
	file, fileheader, err := r.FormFile("uploadFile")
	if err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
	defer file.Close()
	fileSize := fileheader.Size
	if fileSize > maxUploadSize {
		render.Render(w, r, models.ErrorRenderer(errors.New("file too big")))
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg", "image/jpg", "image/png":
		break
	default:
		render.Render(w, r, models.ErrorRenderer(errors.New("invalid file type")))
		return
	}
	newFileName := fmt.Sprintf("%d", time.Now().Unix())
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		render.Render(w, r, models.ErrorRenderer(errors.New("cant read file type")))
		return
	}
	newPath := filepath.Join(uploadPath, newFileName+fileEndings[0])
	newFile, err := os.Create(newPath)
	if err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
	defer newFile.Close()
	if _, err = newFile.Write(fileBytes); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(errors.New("cant write file")))
		return
	}
	dorayaki.ImageURL = Host + FilesURI + newFileName + fileEndings[0]
	if rs := database.DB.Save(&dorayaki); rs.Error != nil {
		render.Render(w, r, models.ErrorRenderer(rs.Error))
		return
	}
	resp := models.ResponseDorayaki{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, dorayaki)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func paginateDorayakiGorm(w http.ResponseWriter, r *http.Request) {
	var data []models.Dorayaki
	rasa := r.URL.Query().Get("dorayaki")
	cond := database.DB.Where("rasa LIKE ?", "%"+rasa+"%")
	totalItems := cond.Find(&models.Dorayaki{}).RowsAffected
	sort := r.URL.Query()["sort"]
	idxPage, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
	if err != nil {
		idxPage = 1
	}
	itemsPerPage, err := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	if err != nil {
		itemsPerPage = 10
	}
	if sort != nil {
		var orderBy string
		for _, sorts := range sort {
			orderBy = sorts
			if sorts[0] == '-' {
				orderBy = sorts[1:] + " desc"
			}
			cond = cond.Order(orderBy)
		}
	}

	cond = cond.Limit(itemsPerPage).Offset(itemsPerPage * (idxPage - 1))
	if rs := cond.Find(&data); rs.Error != nil {
		render.Render(w, r, models.ServerErrorRenderer(rs.Error))
		return
	}
	respPaginate := models.ResponsePaginate{
		Response:     *models.SuccessResponse,
		ItemsPerPage: int64(itemsPerPage),
		TotalItems:   totalItems,
		PageIndex:    int64(idxPage),
		TotalPages:   int64(math.Ceil(float64(totalItems) / float64(itemsPerPage))),
		Sort:         sort,
	}
	resp := models.ResponsePaginateDorayaki{
		ResponsePaginate: respPaginate,
		Rasa:             r.URL.Query().Get("dorayaki"),
		Data:             data,
	}
	if err = render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func DorayakiContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dorayakiId := chi.URLParam(r, "dorayakiId")
		if dorayakiId == "" {
			render.Render(w, r, models.ErrorRenderer(fmt.Errorf("dorayaki ID is required")))
			return
		}
		id, err := strconv.Atoi(dorayakiId)
		if err != nil {
			render.Render(w, r, models.ErrorRenderer(fmt.Errorf("invalid dorayaki ID")))
			return
		}
		ctx := context.WithValue(r.Context(), keyDorayaki, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
