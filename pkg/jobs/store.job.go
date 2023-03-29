package jobs

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/afroborg/apk-rear/pkg/models"
	"github.com/afroborg/apk-rear/pkg/utils"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func SyncStores(DB *gorm.DB, c *cron.Cron) {
	log.Println("Initializing store job")

	var count int64

	DB.Model(&models.Store{}).Count(&count)

	if count == 0 {
		log.Print("No stores found, running job prematurely")
		getStores(DB)
	}

	// Run every first day of the month
	c.AddFunc("0 0 1 1 * *", func() {
		go getStores(DB)
	})
}

func getStores(DB *gorm.DB) {

	stores := fetchStores()

	if len(stores) != 0 {
		DB.Exec("DELETE FROM stores")
	}

	log.Printf("\nFetched %d stores\n-----------------", len(stores))

	if len(stores) != 0 {
		DB.Exec("DELETE FROM stores")
	}

	DB.CreateInBatches(&stores, 500)
}

func fetchStores() []models.Store {

	client := &http.Client{}
	subscriptionKey := utils.GetEnvVariable("SYSTEMBOLAGET_STORE_KEY", "")

	url := "https://api-extern.systembolaget.se/site/v2/store"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header = http.Header{
		"Accept":                    []string{"application/json"},
		"Ocp-Apim-Subscription-Key": []string{subscriptionKey},
	}

	res, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return []models.Store{}
	}

	defer res.Body.Close()

	data := []models.StoreResponse{}
	json.NewDecoder(res.Body).Decode(&data)

	if len(data) == 0 {
		log.Print("Found 0 stores")
		return []models.Store{}
	}

	log.Printf("Fetched %d stores", len(data))

	return mapStores(data)
}

func mapStores(response []models.StoreResponse) []models.Store {
	list := []models.Store{}

	for _, product := range response {
		store := models.Store{
			Name:            product.Alias,
			SystembolagetId: product.SiteId,
			Address:         product.Address,
			PostalCode:      product.PostalCode,
			City:            product.City,
		}

		list = append(list, store)
	}

	return list
}
