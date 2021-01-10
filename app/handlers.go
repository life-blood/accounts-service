package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("GET").
		Path("/").
		HandlerFunc(app.homePage)

	app.Router.
		Methods("GET").
		Path("/accounts/donors").
		HandlerFunc(app.getAllDonors)

	app.Router.
		Methods("GET").
		Path("/accounts/donors/{id:[0-9]+}").
		HandlerFunc(app.getDonorByID)

	app.Router.
		Methods("GET").
		Path("/accounts/acceptors").
		HandlerFunc(app.getAllAcceptors)

	app.Router.
		Methods("PUT").
		Path("/accounts/donors/{id:[0-9]+}").
		HandlerFunc(app.updateDonorByID)

	app.Router.
		Methods("GET").
		Path("/accounts/donors/{bloodGroup:[a-z][A-Z]+}").
		HandlerFunc(app.getDonorsByBloodGroup)

	app.Router.
		Methods("GET").
		Path("/accounts/acceptors/{bloodGroup:[a-z][A-Z]+}").
		HandlerFunc(app.getAcceptorsByBloodGroup)
}

func (app *App) homePage(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	log.Printf("Endpoint Hit: GET /")
}

func (app *App) getAllDonors(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/donors")
}

func (app *App) getAllAcceptors(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/acceptors")
}

func (app *App) getDonorByID(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/donors/:id")
}

func (app *App) updateDonorByID(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Endpoint Hit: PUT /accounts/donors/:id")
}

func (app *App) updateAcceptorByID(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Endpoint Hit: PUT /accounts/acceptors/:id")
}

func (app *App) getDonorsByBloodGroup(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/donors/:bloodGroup")
}

func (app *App) getAcceptorsByBloodGroup(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/acceptors/:bloodGroup")
}
