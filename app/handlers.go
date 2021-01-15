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
	Router   *mux.Router
	Database *sql.DB
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
		Methods("GET").
		Path("/accounts/donors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.getDonorByID)

	app.Router.
		Methods("GET").
		Path("/accounts/acceptors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.getAcceptorByID)

	app.Router.
		Methods("GET").
		Path("/accounts/acceptors").
		HandlerFunc(app.getAllAcceptors)

	app.Router.
		Methods("PUT").
		Path("/accounts/donors/{id:[a-zA-Z0-9]+}").
		HandlerFunc(app.updateDonorByID)

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
	donors := make([]Donor, 0)
	rows, err := app.Database.Query(`SELECT * FROM donors;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, lastname, phone, email, age, gender, bloodGroup, city, regDate string
		err := rows.Scan(&id, &name, &lastname, &phone, &email, &age, &gender, &bloodGroup, &city, &regDate)
		if err != nil {
			log.Fatal(err)
		}

		donors = append(donors, Donor{
			ID:               id,
			FirstName:        name,
			LastName:         lastname,
			PhoneNumber:      phone,
			Email:            email,
			Age:              age,
			Gender:           gender,
			BloodGroup:       bloodGroup,
			City:             city,
			RegistrationDate: regDate})
	}
	if err := rows.Err(); err != nil {
		log.Printf(err.Error())
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(donors); err != nil {
		panic(err)
	}
}

func (app *App) getAllAcceptors(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/acceptors")
	acceptors := make([]Acceptor, 0)
	rows, err := app.Database.Query(`SELECT * FROM acceptors;`)
	switch {
	case err != nil:
		log.Fatal(err)

	case err == sql.ErrNoRows:
		log.Printf("No acceptors found.")
		w.WriteHeader(http.StatusNotFound)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name, lastname, bloodGroup, city, bloodCenter, regDate string
		err := rows.Scan(&id, &name, &lastname, &bloodGroup, &city, &bloodCenter, &regDate)
		if err != nil {
			log.Fatal(err)
		}

		acceptors = append(acceptors, Acceptor{
			ID:               id,
			FirstName:        name,
			LastName:         lastname,
			BloodGroup:       bloodGroup,
			City:             city,
			BloodCenter:      bloodCenter,
			RegistrationDate: regDate})
	}
	if err := rows.Err(); err != nil {
		log.Printf(err.Error())
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(acceptors); err != nil {
		panic(err)
	}
}

func (app *App) getDonorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/donors/:id")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Fatal("No ID in the path")
	}
	donor := Donor{}
	err := app.Database.QueryRow(`SELECT * FROM donors WHERE id=?`, id).Scan(
		&donor.ID,
		&donor.FirstName,
		&donor.LastName,
		&donor.PhoneNumber,
		&donor.Email,
		&donor.Age,
		&donor.Gender,
		&donor.BloodGroup,
		&donor.City,
		&donor.RegistrationDate)
	if err != nil {
		log.Printf(err.Error())
		log.Fatal("Database SELECT failed")
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(donor); err != nil {
		panic(err)
	}
}

func (app *App) getAcceptorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/acceptors/:id")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Fatal("No ID in the path")
	}
	acceptor := Acceptor{}
	err := app.Database.QueryRow(`SELECT * FROM acceptors WHERE id=?`, id).Scan(
		&acceptor.ID,
		&acceptor.FirstName,
		&acceptor.LastName,
		&acceptor.BloodGroup,
		&acceptor.City,
		&acceptor.BloodCenter,
		&acceptor.RegistrationDate)
	if err != nil {
		log.Printf(err.Error())
		log.Fatal("Database SELECT failed")
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(acceptor); err != nil {
		panic(err)
	}
}

func (app *App) updateDonorByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: PUT /accounts/donors/:id")
	// vars := mux.Vars(r)
	// id, ok := vars["id"]
	// if !ok {
	// 	log.Fatal("No ID in the path")
	// }
	// donor := Donor{}
	// err := app.Database.QueryRow(`SELECT * FROM donors WHERE id=?`, id).Scan(
	// 	&donor.ID,
	// 	&donor.FirstName,
	// 	&donor.LastName,
	// 	&donor.PhoneNumber,
	// 	&donor.Email,
	// 	&donor.Age,
	// 	&donor.Gender,
	// 	&donor.BloodGroup,
	// 	&donor.City,
	// 	&donor.RegistrationDate)
	// if err != nil {
	// 	log.Printf(err.Error())
	// 	log.Fatal("Database SELECT failed")
	// }

	// w.WriteHeader(http.StatusOK)
	// if err := json.NewEncoder(w).Encode(donor); err != nil {
	// 	panic(err)
	// }

}

func (app *App) updateAcceptorByID(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Endpoint Hit: PUT /accounts/acceptors/:id")
}

func (app *App) getDonorsByBloodGroup(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/donors/bloodtype/:bloodGroup")
	vars := mux.Vars(r)
	bloodGroup, ok := vars["bloodGroup"]
	if !ok {
		log.Fatal("No bloodGroup in the path")
	}
	donorsWithBloodGroup := make([]Donor, 0)
	rows, err := app.Database.Query(`SELECT * FROM donors WHERE bloodGroup=?`, bloodGroup)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, lastname, phone, email, age, gender, bloodGroup, city, regDate string
		err := rows.Scan(&id, &name, &lastname, &phone, &email, &age, &gender, &bloodGroup, &city, &regDate)
		if err != nil {
			log.Fatal(err)
		}

		donorsWithBloodGroup = append(donorsWithBloodGroup, Donor{
			ID:               id,
			FirstName:        name,
			LastName:         lastname,
			PhoneNumber:      phone,
			Email:            email,
			Age:              age,
			Gender:           gender,
			BloodGroup:       bloodGroup,
			City:             city,
			RegistrationDate: regDate})
	}
	if err := rows.Err(); err != nil {
		log.Printf(err.Error())
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(donorsWithBloodGroup); err != nil {
		panic(err)
	}
}

func (app *App) getAcceptorsByBloodGroup(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: GET /accounts/acceptors/bloodtype/:bloodGroup")
	vars := mux.Vars(r)
	bloodGroup, ok := vars["bloodGroup"]
	if !ok {
		log.Fatal("No bloodGroup in the path")
	}
	acceptorsWithBloodGroup := make([]Acceptor, 0)
	rows, err := app.Database.Query(`SELECT * FROM acceptors WHERE bloodGroup=?`, bloodGroup)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, lastname, bloodGroup, city, bloodCenter, regDate string
		err := rows.Scan(&id, &name, &lastname, &bloodGroup, &city, &bloodCenter, &regDate)
		if err != nil {
			log.Fatal(err)
		}

		acceptorsWithBloodGroup = append(acceptorsWithBloodGroup, Acceptor{
			ID:               id,
			FirstName:        name,
			LastName:         lastname,
			BloodGroup:       bloodGroup,
			City:             city,
			BloodCenter:      bloodCenter,
			RegistrationDate: regDate})
	}
	if err := rows.Err(); err != nil {
		log.Printf(err.Error())
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(acceptorsWithBloodGroup); err != nil {
		panic(err)
	}
}

func (app *App) addDonor(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint Hit: POST /accounts/donors")
	stmt, err := app.Database.Prepare(`INSERT INTO donors (id, name, lastName, phone, email, age, gender, bloodGroup, city, regDate)
	VALUES (?,?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	reqData := make(map[string]string)
	json.Unmarshal(body, &reqData)
	id := shortuuid.New()
	name := reqData["name"]
	lastName := reqData["lastName"]
	phone := reqData["phone"]
	email := reqData["email"]
	age := reqData["age"]
	gender := reqData["gender"]
	bloodGroup := reqData["bloodGroup"]
	city := reqData["city"]
	timeNow := time.Now()
	regDate := timeNow.Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(id, name, lastName, phone, email, age, gender, bloodGroup, city, regDate)
	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		log.Printf(err.Error())
		log.Fatal("Database INSERT to donors failed")
	}
	w.WriteHeader(http.StatusOK)

}
