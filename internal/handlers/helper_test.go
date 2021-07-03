package handlers

import (
	"dorayaki/internal/helpers"
	"net/http"
	"testing"
)

func TestPaginateAbstract(t *testing.T) {
	// assert := assert.New(t)
	r, err := http.NewRequest("POST", "localhost:8080/search?dorayaki=coklat", nil)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	rawString, rawArgs, _ := helpers.PaginateAbstract("dorayaki", r)
	t.Logf("%s\n%v\n", rawString, rawArgs)
}
