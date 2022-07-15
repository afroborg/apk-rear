package models

type StoreResponse struct {
	SiteId     string
	Alias      string
	Address    string
	PostalCode string
	City       string
}

type Store struct {
	ID              int    `gorm:"primaryKey;not null;autoIncrement:true"`
	Name            string `json:"name"`
	SystembolagetId string `json:"systembolagetId" gorm:"primaryKey"`
	Address         string `json:"address"`
	PostalCode      string `json:"postalCode"`
	City            string `json:"city"`
}
