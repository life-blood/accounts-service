package app

import (
	"database/sql"
	"log"
)

//DonorsMySQL mysql repo
type DonorsMySQL struct {
	db *sql.DB
}

//NewDonorsMySQL create new repository
func NewDonorsMySQL(db *sql.DB) *DonorsMySQL {
	return &DonorsMySQL{
		db: db,
	}
}

//Create a Donor
func (r *DonorsMySQL) Create(donor Donor) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO donors (id, name, lastName, phone, email, age, gender, bloodGroup, city, regDate)
		VALUES (?,?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		donor.ID, donor.FirstName, donor.LastName, donor.PhoneNumber, donor.Email, donor.Age, donor.Gender, donor.BloodGroup, donor.City, donor.RegistrationDate,
	)
	return err
}

//GetAll donors
func (r *DonorsMySQL) GetAll() ([]Donor, error) {
	donors := make([]Donor, 0)
	rows, err := r.db.Query(`SELECT * FROM donors;`)
	if err != nil {
		log.Printf(err.Error())
		return donors, err
	}
	defer rows.Close()

	donorsCount := 0
	for rows.Next() {
		var id, name, lastname, phone, email, age, gender, bloodGroup, city, regDate string
		err := rows.Scan(&id, &name, &lastname, &phone, &email, &age, &gender, &bloodGroup, &city, &regDate)
		if err != nil {
			log.Printf(err.Error())
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
		donorsCount++
	}

	return donors, nil
}

//GetByID Retrieve a donor by Id
func (r *DonorsMySQL) GetByID(id string) (Donor, error) {
	donor := Donor{}
	err := r.db.QueryRow(`SELECT * FROM donors WHERE id=?`, id).Scan(
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

	return donor, err
}

//Update donor by ID
func (r *DonorsMySQL) Update(donor Donor) (Donor, error) {
	stmt, err := r.db.Prepare(`UPDATE donors SET name=?,lastName=?,phone=?,email=?,age=?,gender=?,city=? WHERE id=?;`)
	if err != nil {
		log.Printf(err.Error())
	}

	_, err = stmt.Exec(donor.FirstName, donor.LastName, donor.PhoneNumber, donor.Email, donor.Age, donor.Gender, donor.City, donor.ID)
	if err != nil {
		log.Printf(err.Error())
	}

	return donor, err
}

//DeleteByID check whether donor exists and remove
func (r *DonorsMySQL) DeleteByID(id string) error {
	_, err := r.db.Query(`DELETE FROM donors WHERE id=?`, id)
	return err
}
