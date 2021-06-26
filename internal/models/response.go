package models

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseDorayaki struct {
	Response
	Data []Dorayaki `json:"data"`
}

type ResponseToko struct {
	Response
	Data []Toko `json:"data"`
}
