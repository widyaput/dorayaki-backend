package helpers

import (
	"dorayaki/internal/models"
	"net/http"
	"strconv"

	sq "github.com/Masterminds/squirrel"
)

func PaginateAbstract(tableName string, r *http.Request) (string, []interface{}, error) {
	sql := TakeQuery(tableName, r)
	sort := r.URL.Query().Get("sort")
	pageIndex := r.URL.Query().Get("pageIndex")
	itemsPerPage := r.URL.Query().Get("itemsPerPage")
	orderBy := "updated_at"
	if sort != "" {
		orderBy = sort
		if sort[0] == '-' {
			orderBy = sort[1:] + " desc"
		}
	}
	sql = sql.OrderBy(orderBy)
	var idxPage int
	var itemPage int
	idxPage, err := strconv.Atoi(pageIndex)
	if err != nil {
		idxPage = 1
	}
	itemPage, err = strconv.Atoi(itemsPerPage)
	if err != nil {
		itemPage = 10
	}
	sql = sql.Offset(uint64(itemPage) * (uint64(idxPage) - 1))
	sql = sql.Limit(uint64(itemPage))
	resultQuery, resultArgs, err := sql.ToSql()
	return resultQuery, resultArgs, err
}

func TakeQuery(tableName string, r *http.Request) sq.SelectBuilder {
	sql := sq.Select("*").From(tableName)

	var dorayaki, kecamatan, provinsi string
	query := sq.And{}
	if tableName == models.Dorayaki.TableName(models.Dorayaki{}) {
		dorayaki = r.URL.Query().Get("dorayaki")
		query = append(query, sq.Like{"rasa": "%" + dorayaki + "%"})
	}
	if tableName == models.Toko.TableName(models.Toko{}) {
		kecamatan = r.URL.Query().Get("kecamatan")
		provinsi = r.URL.Query().Get("provinsi")
		if kecamatan != "" {
			query = append(query, sq.Like{"kecamatan": "%" + kecamatan + "%"})
		}
		if provinsi != "" {
			query = append(query, sq.Like{"provinsi": "%" + provinsi + "%"})
		}
	}
	sql = sql.Where(query)
	return sql
}
