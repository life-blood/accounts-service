package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
)

// App is a wrapper struct over Router and Database used to manage db connections and endpoint requests centrally
type App struct {
	Router        *mux.Router
	DonorsRepo    *DonorsMySQL
	AcceptorsRepo *AcceptorsMySQL
}

// SetupRouter is used to provide mapping between different endpoints hit and handler functions
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
		Methods("POST").
		Path("/accounts/donors").
		HandlerFunc(app.addDonor)

	app.Router.
		Methods("POST").
		Path("/accounts/acceptors").
		HandlerFunc(app.addAcceptor)

	app.Router.
		Methods("GET").
		Path("/accounts/donors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.getDonorByID)

	app.Router.
		Methods("DELETE", "OPTIONS").
		Path("/accounts/donors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.deleteDonorByID)

	app.Router.
		Methods("GET").
		Path("/accounts/acceptors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.getAcceptorByID)

	app.Router.
		Methods("DELETE", "OPTIONS").
		Path("/accounts/acceptors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.deleteAcceptorByID)

	app.Router.
		Methods("GET").
		Path("/accounts/acceptors").
		HandlerFunc(app.getAllAcceptors)

	app.Router.
		Methods("PUT", "OPTIONS").
		Path("/accounts/donors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.updateDonorByID)

	app.Router.
		Methods("PUT", "OPTIONS").
		Path("/accounts/acceptors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.updateAcceptorByID)

	app.Router.
		Methods("GET").
		Path("/accounts/donors/bloodtype/{bloodGroup:[a-zA-Z0-9]+}").
		HandlerFunc(app.getDonorsByBloodGroup)

	app.Router.
		Methods("GET").
		Path("/accounts/acceptors/bloodtype/{bloodGroup:[a-zA-Z0-9]+}").
		HandlerFunc(app.getAcceptorsByBloodGroup)
}

func (app *App) homePage(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	log.Printf("Endpoint Hit: GET /")
}

func (app *App) getAllDonors(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/donors")

	setupCORS(&w, r)

	donors, err := app.DonorsRepo.GetAll()

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(donors); err != nil {
			log.Printf(err.Error())
		}
	}
}

func (app *App) getAllAcceptors(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/acceptors")
	setupCORS(&w, r)

	acceptors, err := app.AcceptorsRepo.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Could not load acceptors.")
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(acceptors); err != nil {
			log.Printf(err.Error())
		}
	}
}

func (app *App) getDonorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/donors/:id")
	setupCORS(&w, r)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Printf("No ID in the path for GET /accounts/donors/:id")
	}

	donor, err := app.DonorsRepo.GetByID(id)

	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(donor); err != nil {
		log.Printf(err.Error())
	}
}

func (app *App) getAcceptorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/acceptors/:id")
	setupCORS(&w, r)

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Printf("No ID in the path for GET /accounts/acceptors/:id")
	}

	acceptor, err := app.AcceptorsRepo.GetByID(id)

	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Database SELECT from acceptors failed for GET /accounts/acceptors/:id")
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(acceptor); err != nil {
			log.Printf(err.Error())
		}
	}
}

func (app *App) updateDonorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: PUT /accounts/donors/:id")

	setupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Printf("No ID in the path for PUT /accounts/donors/:id")
	}

	donor, err := app.DonorsRepo.GetByID(id)

	if err == sql.ErrNoRows {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Cannot update unexisting donor.")
	}
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Database SELECT from donors failed for PUT /accounts/donors/:id")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf(err.Error())
	}
	reqData := make(map[string]string)
	json.Unmarshal(body, &reqData)

	if firstName, exists := reqData["name"]; exists {
		donor.FirstName = firstName
	}
	if lastName, exists := reqData["lastName"]; exists {
		donor.LastName = lastName
	}
	if phone, exists := reqData["phone"]; exists {
		donor.PhoneNumber = phone
	}
	if email, exists := reqData["email"]; exists {
		donor.Email = email
	}
	if age, exists := reqData["age"]; exists {
		donor.Age = age
	}
	if gender, exists := reqData["gender"]; exists {
		donor.Gender = gender
	}
	if city, exists := reqData["city"]; exists {
		donor.City = city
	}

	_, err = app.DonorsRepo.Update(donor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *App) updateAcceptorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: PUT /accounts/acceptors/:id")
	setupCORS(&w, r)

	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Printf("No ID in the path PUT /accounts/acceptors/:id")
	}
	acceptor, err := app.AcceptorsRepo.GetByID(id)

	if err != nil {
		log.Printf(err.Error())
		log.Printf("Cannot update unexisting acceptor.")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf(err.Error())
	}

	reqData := make(map[string]string)
	json.Unmarshal(body, &reqData)

	if firstName, exists := reqData["name"]; exists {
		acceptor.FirstName = firstName
	}
	if lastName, exists := reqData["lastName"]; exists {
		acceptor.LastName = lastName
	}
	if city, exists := reqData["city"]; exists {
		acceptor.City = city
	}
	if bloodCenter, exists := reqData["bloodCenter"]; exists {
		acceptor.BloodCenter = bloodCenter
	}

	err = app.AcceptorsRepo.Update(acceptor)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *App) getDonorsByBloodGroup(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/donors/bloodtype/:bloodGroup")
	setupCORS(&w, r)
	vars := mux.Vars(r)
	bloodGroup, ok := vars["bloodGroup"]
	if !ok {
		log.Printf("No bloodGroup in the path for GET /accounts/donors/bloodtype/:bloodGroup")
	}

	donors, err := app.DonorsRepo.GetByBloodGroup(bloodGroup)

	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(donors); err != nil {
			log.Printf(err.Error())
		}
	}
}

func (app *App) getAcceptorsByBloodGroup(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/acceptors/bloodtype/:bloodGroup")
	setupCORS(&w, r)
	vars := mux.Vars(r)
	bloodGroup, ok := vars["bloodGroup"]
	if !ok {
		log.Printf("No bloodGroup in the path for GET /accounts/acceptors/bloodtype/:bloodGroup")
	}

	acceptors, err := app.AcceptorsRepo.GetByBloodGroup(bloodGroup)

	if err != nil {
		log.Printf("No donors with the specific bloodGroup found.")
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(acceptors); err != nil {
			log.Printf(err.Error())
		}
	}
}

func (app *App) addDonor(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: POST /accounts/donors")
	setupCORS(&w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf(err.Error())
	}
	reqData := make(map[string]string)
	json.Unmarshal(body, &reqData)
	donor := Donor{}
	donor.ID = shortuuid.New()
	donor.FirstName = reqData["name"]
	donor.LastName = reqData["lastName"]
	donor.BloodGroup = reqData["bloodGroup"]
	donor.City = reqData["city"]
	donor.PhoneNumber = reqData["phone"]
	donor.Email = reqData["email"]
	donor.Age = reqData["age"]
	donor.Gender = reqData["gender"]
	donor.City = reqData["city"]
	timeNow := time.Now()
	donor.RegistrationDate = timeNow.Format("2006-01-02 15:04:05")

	app.DonorsRepo.Create(donor)

	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Database INSERT to donors failed")
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (app *App) addAcceptor(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: POST /accounts/acceptors")
	setupCORS(&w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf(err.Error())
	}
	reqData := make(map[string]string)
	json.Unmarshal(body, &reqData)
	acceptor := Acceptor{}
	acceptor.ID = shortuuid.New()
	acceptor.FirstName = reqData["name"]
	acceptor.LastName = reqData["lastName"]
	acceptor.BloodCenter = reqData["bloodCenter"]
	acceptor.BloodGroup = reqData["bloodGroup"]
	acceptor.City = reqData["city"]
	timeNow := time.Now()
	acceptor.RegistrationDate = timeNow.Format("2006-01-02 15:04:05")

	err = app.AcceptorsRepo.Create(acceptor)

	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Database INSERT to acceptors failed")
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (app *App) deleteDonorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: DELETE /accounts/donors/:id")

	setupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Printf("No ID in the path for DELETE /accounts/donors/:id")
	}

	err := app.DonorsRepo.DeleteByID(id)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Database DELETE failed for DELETE /accounts/donors/:id")
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("X"))
	}
}

func (app *App) deleteAcceptorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: DELETE /accounts/acceptor/:id")

	setupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Printf("No ID in the path for DELETE /accounts/acceptor/:id")
	}
	err := app.AcceptorsRepo.DeleteByID(id)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Database DELETE failed for DELETE /accounts/acceptor/:id")
	}

	w.WriteHeader(http.StatusOK)
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
