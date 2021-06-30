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

func (re *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, re.StatusCode)
	return nil
}
