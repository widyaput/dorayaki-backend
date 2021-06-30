package handlers

import (
	"dorayaki/internal/models"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const maxUploadSize = 2 * 1024 * 1024
const uploadPath = "internal/assets/"

// const relativeUploadPath = "../assets/"

func assets(router chi.Router) {
	router.Post("/upload", uploadImage)
	fs := http.FileServer(http.Dir(uploadPath))
	// fmt.Println(http.Dir(uploadPath))
	router.Handle("/files/*", http.StripPrefix("/assets/files/", fs))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
	// fileType := r.PostFormValue("type")
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
	newFileName := randToken(12)
	newFileName = fmt.Sprintf("%d", time.Now()) + newFileName
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		render.Render(w, r, models.ErrorRenderer(errors.New("cant read file type")))
		return
	}
	newPath := filepath.Join(uploadPath, newFileName+fileEndings[0])
	newFile, err := os.Create(newPath)
	if err != nil {
		render.Render(w, r, models.ErrorRenderer(err))
		return
	}
	defer newFile.Close()
	if _, err = newFile.Write(fileBytes); err != nil {
		render.Render(w, r, models.ErrorRenderer(errors.New("cant write file")))
		return
	}
	// resp := models.ResponseImage{Response: *models.SuccessResponse}
	// resp.Data = append(resp.Data, newPath)
	// if err = render.Render(w, r, &resp); err != nil {
	// 	render.Render(w, r, models.ServerErrorRenderer(err))
	// 	return
	// }
}
