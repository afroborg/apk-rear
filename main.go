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
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func setupRouter(DB *gorm.DB) error {
	r := mux.NewRouter()

	h := handlers.New(DB)

	// Set up middleware logger
	r.Use(middlewares.Logger)

	r.Path("/api/alcohols").Methods("GET").Queries("page", "{page:[0-9]+}", "per_page", "{per_page:[0-9]+}").HandlerFunc(h.GetAlcohols)
	r.Path("/api/stores").Methods("GET").HandlerFunc(h.GetStores)
	r.HandleFunc("/api/health", h.Health).Methods("GET", "POST", "PUT")
	r.HandleFunc("/api/status", h.Status).Methods("GET")

	port := utils.GetEnvVariable("PORT", "8080")

	log.Println("Starting server on port " + port)

	err := http.ListenAndServe(":"+port, r)

	return err
}

func main() {
	godotenv.Load()
	log.Println("Intializing server...")

	DB := db.Init()
	c := cron.New()

	jobs.SyncStores(DB, c)
	jobs.SyncAlcohol(DB, c)
	c.Start()

	printCronEntries(c.Entries())

	log.Fatalf("Failed to setup router, %v", setupRouter(DB))
}

func printCronEntries(cronEntries []*cron.Entry) {
	log.Printf("Cron Info: %+v\n", cronEntries)
}
