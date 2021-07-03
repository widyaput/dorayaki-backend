package helpers

import (
	"dorayaki/internal/models"
	"net/http"
	"strconv"

	sq "github.com/Masterminds/squirrel"
)

// TODO:Paginate.
// Paginate idea : take dorayaki based rasa and shows some shops that sell it. Sorting based on freshly updated, fresh category, fresh shop.
func PaginateAbstract(tableName string, r *http.Request) (string, []interface{}, error) {
	var sql sq.SelectBuilder
	sort := r.URL.Query().Get("sort")
	pageIndex := r.URL.Query().Get("pageIndex")
	itemsPerPage := r.URL.Query().Get("itemsPerPage")
	var dorayaki, kecamatan, provinsi string
	query := sq.And{}
	if tableName == models.Dorayaki.TableName(models.Dorayaki{}) {
		// TODO: make this effective
		sql = sq.Select("*").From(tableName)
		dorayaki = r.URL.Query().Get("dorayaki")
		sql = sql.Join(models.Toko{}.TableName()).Join(models.TokoDorayaki{}.TableName()).
			Where(sq.Eq{
				models.TokoDorayaki{}.TableName() + ".dorayaki_id": "id",
				models.Toko{}.TableName() + ".id":                  models.TokoDorayaki{}.TableName() + ".toko_id"})
		query = append(query, sq.Like{"rasa": "%" + dorayaki + "%"})
	}
	if tableName == models.Toko.TableName(models.Toko{}) {
		sql = sq.Select("*").From(tableName)
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
	if sort != "" {
		orderBy := sort
		if sort[0] == '-' {
			orderBy = sort[1:] + " desc"
		}
		sql = sql.OrderBy(orderBy)
	}
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
