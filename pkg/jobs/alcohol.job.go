package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/afroborg/apk-rear/pkg/models"
	"github.com/afroborg/apk-rear/pkg/utils"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func SyncAlcohol(DB *gorm.DB, c *cron.Cron) {
	log.Println("Initializing alcohol job")

	var count int64

	DB.Model(&models.Alcohol{}).Count(&count)

	if count == 0 {
		log.Print("No alcohols found, running job prematurely")
		go run(DB)
	}

	// Run every sunday at 03:00
	c.AddFunc("0 0 3 * * 0", func() {
		go run(DB)
	})
}

func run(DB *gorm.DB) {

	alcohols := fetchAlcohols()

	log.Printf("\nFetched %d alcohols\n-----------------", len(alcohols))

	if len(alcohols) != 0 {
		DB.Exec("DELETE FROM alcohols")
	}

	DB.CreateInBatches(&alcohols, 500)
}

func fetchAlcohols() []models.Alcohol {
	alcohols := []models.Alcohol{}
	page := 1

	client := &http.Client{}
	subscriptionKey := utils.GetEnvVariable("SYSTEMBOLAGET_SUBSCRIPTION_KEY", "")

	for {
		url := "https://api-extern.systembolaget.se/sb-api-ecommerce/v1/productsearch/search?size=30&page=" + fmt.Sprint(page)
		req, _ := http.NewRequest("GET", url, nil)

		req.Header = http.Header{
			"Accept":                    []string{"application/json"},
			"Ocp-Apim-Subscription-Key": []string{subscriptionKey},
		}

		res, err := client.Do(req)

		if err != nil {
			log.Println(err)
			break
		}

		defer res.Body.Close()

		data := models.AlcoholResponse{}
		json.NewDecoder(res.Body).Decode(&data)

		if len(data.Products) == 0 {
			log.Printf("Found 0 alochols on page %d", page)
			break
		}

		log.Printf("Fetched %d alcohols from page %d, next page: %d", len(data.Products), page, data.Metadata.NextPage)

		alcohols = append(alcohols, mapAlcohols(data)...)

		page = data.Metadata.NextPage
	}

	return alcohols

}

func mapAlcohols(response models.AlcoholResponse) []models.Alcohol {
	list := []models.Alcohol{}

	for _, product := range response.Products {
		alcohol := models.Alcohol{
			Name:            product.ProductNameBold,
			SystembolagetId: product.ProductId,
			Number:          product.ProductNumber,
			Price:           product.Price,
			Volume:          product.Volume,
			Country:         product.Country,
			Image:           getImage(product.Images),
			Apk:             calcApk(product.AlcoholPercentage, product.Volume, product.Price),
			Category:        product.CategoryLevel1,
			Link: 			 "https://systembolaget.se/" + product.ProductNumber,
		}

		list = append(list, alcohol)
	}

	return list
}

func getImage(images []models.AlcoholResponseImages) string {
	if len(images) < 1 {
		return ""
	}

	return images[0].ImageUrl
}

func calcApk(percentage float64, volume float64, price float64) float64 {
	return (percentage / 100) * volume / price
}
