package main

import (
    "fmt"
    "log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

const port = 4200

type Donor struct {
    Id string `json:id`
    FirstName string `json:name`
	LastName string `json:lastName`
	PhoneNumber string `json:phone`
	Email string `json:email`
	Age string `json:age`
	Gender string `json:gender`
	BloodGroup string `json:bloodGroup`
	Adress string `json:address`
	HealthStatus string `json:healthStatus`
	RegistrationDate string `json:regDate`
	LastDonationDate string `json:lastDonationDate`
}

var Donors []Donor

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: GET /")
}

func getAllDonors(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: GET /accounts/donors")
    json.NewEncoder(w).Encode(Donors)
}

func getDonorById(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: GET /accounts/donors/:id")
}

func handleRequests() {
	accountsRouter := mux.NewRouter().StrictSlash(true)
	accountsRouter.HandleFunc("/", homePage)
	accountsRouter.HandleFunc("/accounts/donors", getAllDonors)
	accountsRouter.HandleFunc("/accounts/donors/:id", getDonorById)
	http.Handle("/", accountsRouter)
    log.Fatal(http.ListenAndServe(":4200", accountsRouter))
}

func main() {
	Donors = []Donor{
		Donor{Id:"1", FirstName:"Ivan", LastName:"Petrov", PhoneNumber:"0897656780", Email:"ivan@mail.bg", Age:"31", Gender:"Male", BloodGroup:"AB+", Adress:"Kardzhali", HealthStatus:"Healthy", RegistrationDate:"Sat Dec 12 17:53:21 EET 2010", LastDonationDate:"Sun Mar 15 02:44:15 EET 2019"},
		Donor{Id:"2", FirstName:"Ivanka", LastName:"Petrova", PhoneNumber:"0897656780", Email:"ivanka@mail.bg", Age:"21", Gender:"Female", BloodGroup:"A-", Adress:"Sofia", HealthStatus:"Unhealthy", RegistrationDate:"Sat Nov 12 19:23:25 EET 2012", LastDonationDate:"Sun Mar 15 02:44:15 EET 2020"},
	}

	log.Printf("Starting accounts microservice on port %d", port)
	log.Printf("Rest API v2.0 - Mux Routers")

    handleRequests()
}