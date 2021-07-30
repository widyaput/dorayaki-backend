package handlers

import (
	"dorayaki/configs/database"
	"dorayaki/internal/helpers"
	"dorayaki/internal/models"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const FilesURI = "/api/v1/files/"
const PrefixAPIV1 = "/api/v1/"

var Host = os.Getenv("HOST")

func apiv1(router chi.Router) {
	router.Post("/signin", signIn)
	router.Post("/signout", signOut)
	router.Group(func(router chi.Router) {
		router.Use(authenticator)
		router.Get("/checkProfile", checkAuth)
		router.Post("/uploads", uploadImage)
		router.Delete("/images/{nameOfFile}", deleteImage)
	})
	fs := http.FileServer(http.Dir(uploadPath))
	router.Handle("/files/{nameOfFile}", http.StripPrefix("/api/v1/files/", fs))
	router.Route("/shops", shops)
	router.Route("/dorayakis", dorayakis)
}

func checkAuth(w http.ResponseWriter, r *http.Request) {
	token := findTokens(r, tokenFromHeader, tokenFromCookie, tokenFromQuery)
	// at this moment, token must be available, because we using authenticator middleware
	users := strings.Split(token, ".")
	userDec, _ := b64.StdEncoding.DecodeString(users[1] + "==")
	resp := models.ResponseAuth{Response: *models.SuccessResponse}
	resp.Data.Token = token
	data := struct {
		Username string `json:"username"`
		Exp      int64  `json:"exp"`
	}{}
	json.Unmarshal(userDec, &data)
	resp.Data.Username = data.Username
	resp.Data.Exp = data.Exp
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "nameOfFile")
	if _, err := os.Stat(uploadPath + name); err == nil {
		os.Remove(uploadPath + name)
		return
	}
	render.Render(w, r, models.ErrNotFound)
}

// uploadImage of dorayaki into filesystem
func uploadImage(w http.ResponseWriter, r *http.Request) {
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
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.Mkdir(uploadPath, os.ModePerm)
	}

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
	imageURL := Host + FilesURI + newFileName + fileEndings[0]

	resp := models.ResponseString{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, imageURL)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func signOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
	})
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var cred models.Credentials
	var trueCred models.Credentials
	if err := render.Bind(r, &cred); err != nil {
		render.Render(w, r, models.ErrBadRequest)
		return
	}
	if rs := database.DB.Where("username = ?", cred.Username).First(&trueCred); rs.Error != nil {
		render.Render(w, r, models.ErrUnauthorized)
		return
	}
	if !helpers.CheckPasswordHash(cred.Password, trueCred.Password) {
		render.Render(w, r, models.ErrUnauthorized)
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &models.JWT{
		Username: cred.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	sign := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	token, err := sign.SignedString([]byte(models.JwtKEY))
	if err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
	})
	resp := models.ResponseString{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, token)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func auth(r *http.Request, fn ...func(r *http.Request) string) error {
	token := findTokens(r, fn...)
	if token == "" {
		return fmt.Errorf("no token found")
	}
	return parseToken(token)
}

func findTokens(r *http.Request, fn ...func(r *http.Request) string) string {
	var token string
	for _, funct := range fn {
		token = funct(r)
		if token != "" {
			break
		}
	}
	return token
}

func parseToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWT{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(models.JwtKEY), nil
	})

	if err != nil {
		return err
	}
	if token == nil || !token.Valid {
		return fmt.Errorf("unauthorized")
	}
	return nil
}

func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func tokenFromCookie(r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return c.Value
}
func tokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get("jwt")
}
