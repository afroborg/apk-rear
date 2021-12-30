package main

import (
	"log"
	"net/http"

	"github.com/afroborg/apk-rear/pkg/db"
	"github.com/afroborg/apk-rear/pkg/handlers"
	"github.com/afroborg/apk-rear/pkg/jobs"
	"github.com/afroborg/apk-rear/pkg/middlewares"
	"github.com/afroborg/apk-rear/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func setupRouter(DB *gorm.DB) error {
	r := mux.NewRouter()

	h := handlers.New(DB)

	// Set up middleware logger
	r.Use(middlewares.Logger)

	r.HandleFunc("/alcohols", h.GetAlcohols).Methods("GET")
	r.HandleFunc("/health", h.Health).Methods("GET", "POST", "PUT")

	port := utils.GetEnvVariable("PORT", "8080")

	err := http.ListenAndServe(":"+port, r)

	return err
}

func main() {
	log.Println("Intializing server...")

	DB := db.Init()
	c := cron.New()

	jobs.SyncAlcohol(DB, c)
	c.Start()

	printCronEntries(c.Entries())

	log.Fatalf("Failed to setup router, %v", setupRouter(DB))
}

func printCronEntries(cronEntries []*cron.Entry) {
	log.Printf("Cron Info: %+v\n", cronEntries)
}
