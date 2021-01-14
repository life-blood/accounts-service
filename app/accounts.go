package app

// Donor is a struct used to represent the first account type in LifeBlood system - blood donors
type Donor struct {
	ID               string `json:"id"`
	FirstName        string `json:"name"`
	LastName         string `json:"lastName"`
	PhoneNumber      string `json:"phone"`
	Email            string `json:"email"`
	Age              string `json:"age"`
	Gender           string `json:"gender"`
	BloodGroup       string `json:"bloodGroup"`
	City             string `json:"city"`
	BloodCenter      string `json:"bloodCenter"`
	RegistrationDate string `json:"regDate"`
}

// Acceptor is a struct used to represent the second account type in LifeBlood system - blood acceptors
type Acceptor struct {
	ID               string `json:"id"`
	FirstName        string `json:"name"`
	LastName         string `json:"lastName"`
	BloodGroup       string `json:"bloodGroup"`
	City             string `json:"city"`
	BloodCenter      string `json:"bloodCenter"`
	RegistrationDate string `json:"regDate"`
}
