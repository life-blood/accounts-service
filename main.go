package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/life-blood/accounts-service/app"
	db "github.com/life-blood/accounts-service/config"
)

const port = 4200

func main() {
	//load .env file from given path
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database, err := db.CreateDatabaseConn()
	if err != nil {
		log.Fatalf("Database connection failed: %s", err.Error())
	}

	db.InitializeDatabase(database)
	db.PopulateWithMockData(database)

	donorsRepo := app.NewDonorsMySQL(database)
	acceptorsRepo := app.NewAcceptorsMySQL(database)

	app := &app.App{
		Router:        mux.NewRouter().StrictSlash(true),
		DonorsRepo:    donorsRepo,
		AcceptorsRepo: acceptorsRepo,
	}

	app.SetupRouter()
	log.Printf("Starting accounts microservice on port %d", port)
	log.Fatal(http.ListenAndServe(":4200", app.Router))
}
