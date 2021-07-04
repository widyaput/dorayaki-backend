package models

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

var (
	SuccessResponse       = &Response{StatusCode: 200, Message: "Success"}
	SuccessCreateResponse = &Response{StatusCode: 201, Message: "Create Success"}
)

type ResponseDorayaki struct {
	Response
	Data []Dorayaki `json:"data"`
}

type ResponseToko struct {
	Response
	Data []Toko `json:"data"`
}

type ResponseStok struct {
	Response
	Data []TokoDorayaki
}

type ResponseImageURL struct {
	Response
	Data []string
}

type ResponsePaginate struct {
	Response
	ItemsPerPage int64  `json:"items_per_page"`
	TotalItems   int64  `json:"total_items"`
	PageIndex    int64  `json:"page_index"`
	TotalPages   int64  `json:"total_pages"`
	Sort         string `json:"sort"`
}

type ResponsePaginateToko struct {
	ResponsePaginate
	Kecamatan string `json:"kecamatan"`
	Provinsi  string `json:"provinsi"`
	Data      []Toko `json:"data"`
}

type ResponsePaginateDorayaki struct {
	ResponsePaginate
	Rasa string     `json:"rasa"`
	Data []Dorayaki `json:"data"`
}

func (re *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, re.StatusCode)
	return nil
}
