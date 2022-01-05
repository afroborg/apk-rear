package handlers

import (
	"math"
	"net/http"

	"github.com/afroborg/apk-rear/pkg/models"
	"github.com/afroborg/apk-rear/pkg/utils"
)

func (h handler) GetAlcohols(w http.ResponseWriter, r *http.Request) {
	var alcohols []models.Alcohol
	var totalCount int64
	var fiterCount int64

	q := r.URL.Query()

	// Get query variables
	page := utils.GetQueryVariable(&q, "page", 1).(int)
	perPage := utils.GetQueryVariable(&q, "per_page", 100).(int)
	category := utils.GetQueryVariable(&q, "category", "").(string)
	search := utils.GetQueryVariable(&q, "search", "").(string)

	// Create basic query
	query := h.DB.Model(alcohols).Count(&totalCount).Offset(perPage * (page - 1)).Limit(perPage).Order("apk desc")

	whereStr := ""
	searchArr := make([]interface{}, 0)

	if category != "" {
		whereStr += "category = ?"
		searchArr = append(searchArr, category)
	}

	if search != "" {
		if whereStr != "" {
			whereStr += " AND "
		}

		whereStr += "name ILIKE ? OR thin_name ILIKE ?"
		searchArr = append(searchArr, "%"+search+"%")
		searchArr = append(searchArr, "%"+search+"%")
	}

	if whereStr != "" {
		query = query.Where(whereStr, searchArr...)
	}

	query.Count(&fiterCount).Find(&alcohols)

	h.respond(w, map[string]interface{}{
		"meta": generateMeta(totalCount, fiterCount, page, perPage, len(alcohols)),
		"data": alcohols,
	}, http.StatusOK)
}

func generateMeta(totalCount int64, filterCount int64, page int, perPage int, thisPage int) map[string]interface{} {

	totalPages := int(math.Ceil(float64(filterCount) / float64(perPage)))

	nextPage := page + 1

	if page == totalPages {
		nextPage = -1
	}

	return map[string]interface{}{
		"total":       totalCount,
		"filterCount": filterCount,
		"totalPages":  totalPages,
		"page":        page,
		"nextPage":    nextPage,
		"perPage":     perPage,
		"onThisPage":  thisPage,
	}
}
