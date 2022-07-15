package handlers

import (
	"net/http"

	"github.com/afroborg/apk-rear/pkg/models"
)

func (h handler) GetStores(w http.ResponseWriter, r *http.Request) {
	var stores []models.Store

	h.DB.Model(stores).Find(&stores)

	h.respond(w, map[string]interface{}{
		"data": stores,
	}, http.StatusOK)
}
