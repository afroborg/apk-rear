package handlers

import (
	"math"
	"net/http"

	"github.com/afroborg/apk-rear/pkg/models"
	"github.com/afroborg/apk-rear/pkg/utils"
)

func (h handler) GetAlcohols(w http.ResponseWriter, r *http.Request) {
	var alcohols []models.Alcohol
	var totalAlcohols int64

	q := r.URL.Query()

	// Get query variables
	page := utils.GetQueryVariable(&q, "page", 1).(int)
	perPage := utils.GetQueryVariable(&q, "per_page", 100).(int)
	category := utils.GetQueryVariable(&q, "category", "").(string)
	search := utils.GetQueryVariable(&q, "search", "").(string)

	// Create basic query
	query := h.DB.Model(alcohols).Count(&totalAlcohols).Offset(perPage * (page - 1)).Limit(perPage).Order("apk desc")

	whereStr := ""
	searchArr := make([]interface{}, 0)

	if category != "" {
		whereStr += "Category = ?"
		searchArr = append(searchArr, category)
	}

	if search != "" {
		if whereStr != "" {
			whereStr += " AND "
		}
		whereStr += "Name ILIKE ?"
		searchArr = append(searchArr, "%"+search+"%")
	}

	if whereStr != "" {
		query = query.Where(whereStr, searchArr...)
	}

	query.Find(&alcohols)

	h.respond(w, map[string]interface{}{
		"meta": generateMeta(totalAlcohols, page, perPage, len(alcohols)),
		"data": alcohols,
	}, http.StatusOK)
}

func generateMeta(total int64, page int, perPage int, thisPage int) map[string]interface{} {

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	nextPage := page + 1

	if page == totalPages {
		nextPage = -1
	}

	return map[string]interface{}{
		"total":      total,
		"totalPages": totalPages,
		"page":       page,
		"nextPage":   nextPage,
		"perPage":    perPage,
		"onThisPage": thisPage,
	}
}
