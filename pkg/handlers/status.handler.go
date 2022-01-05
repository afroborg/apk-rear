package handlers

import (
	"net/http"

	"github.com/afroborg/apk-rear/pkg/models"
)

func (h handler) Status(w http.ResponseWriter, r *http.Request) {

	var status = models.Status{}
	h.DB.Find(&models.Status{}).Order("time desc").First(&status)

	h.respond(w, map[string]interface{}{
		"lastSync":       status.Time,
		"syncedProducts": status.SyncedProducts,
	}, 200)
}
