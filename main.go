package main

import (
	"os"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/life-blood/accounts-service/app"
	db "github.com/life-blood/accounts-service/config"
)

func main() {
	database, err := db.CreateDatabaseConn()
	if err != nil {
		log.Fatal("Database connection failed: %s", err.Error())
	}

	port := os.Getenv("PORT")

	app := &app.App{
		Router:   mux.NewRouter().StrictSlash(true),
		Database: database,
	}

	app.SetupRouter()
	log.Printf("Starting accounts microservice on port %d", port)
	log.Fatal(http.ListenAndServe(":" + port, app.Router))
}
