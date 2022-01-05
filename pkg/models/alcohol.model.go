package models

type AlcoholResponseImages struct {
	ImageUrl string
}

type AlcoholResponse struct {
	Metadata struct {
		NextPage int
	}
	Products []struct {
		ProductId         string
		ProductNumber     string
		ProductNameBold   string
		ProductNameThin   string
		AlcoholPercentage float64
		Price             float64
		Volume            float64
		Country           string
		CategoryLevel1    string
		Images            []AlcoholResponseImages
	}
}

type Alcohol struct {
	ID              int     `gorm:"primaryKey;not null;autoIncrement:true"`
	Name            string  `json:"name"`
	ThinName        string  `json:"thinName"`
	SystembolagetId string  `json:"systembolagetId" gorm:"primaryKey"`
	Number          string  `json:"systembolagetNumber" gorm:"primaryKey"`
	Price           float64 `json:"price"`
	Volume          float64 `json:"volume"`
	Country         string  `json:"country"`
	Image           string  `json:"image"`
	Apk             float64 `json:"apk"`
	Category        string  `json:"category"`
	Link            string  `json:"link"`
}
