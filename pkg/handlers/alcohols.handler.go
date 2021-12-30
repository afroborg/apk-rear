package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/afroborg/apk-rear/pkg/models"
)

func (h handler) GetAlcohols(w http.ResponseWriter, r *http.Request) {
	var alcohols []models.Alcohol
	var total int64

	q := r.URL.Query()

	page, err := strconv.Atoi(q.Get("page"))

	if err != nil {
		page = 1
	}

	perPage, err := strconv.Atoi(q.Get("per_page"))

	if err != nil {
		perPage = 100
	}

	h.DB.Model(alcohols).Count(&total).Offset(perPage * (page - 1)).Limit(perPage).Order("apk desc").Find(&alcohols)

	h.respond(w, map[string]interface{}{
		"meta": generateMeta(total, page, perPage, len(alcohols)),
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
