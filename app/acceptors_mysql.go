package app

import (
	"database/sql"
	"log"
)

//AcceptorsMySQL mysql repo
type AcceptorsMySQL struct {
	db *sql.DB
}

//NewAcceptorsMySQL create new repository
func NewAcceptorsMySQL(db *sql.DB) *AcceptorsMySQL {
	return &AcceptorsMySQL{
		db: db,
	}
}

//Create new acceptor
func (r *AcceptorsMySQL) Create(acceptor Acceptor) error {
	stmt, err := r.db.Prepare(
		`INSERT INTO acceptors (id, name, lastName, bloodGroup, city, bloodCenter, regDate)
		VALUES (?,?,?,?,?,?,?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(acceptor.ID, acceptor.FirstName, acceptor.LastName, acceptor.BloodGroup, acceptor.City, acceptor.BloodCenter, acceptor.RegistrationDate)
	return err
}

//GetAll acceptors
func (r *AcceptorsMySQL) GetAll() ([]Acceptor, error) {
	acceptors := make([]Acceptor, 0)
	rows, err := r.db.Query(`SELECT * FROM acceptors;`)
	if err != nil {
		log.Printf(err.Error())
	}
	defer rows.Close()

	acceptorsCount := 0
	for rows.Next() {
		var id, name, lastname, bloodGroup, city, bloodCenter, regDate string
		err := rows.Scan(&id, &name, &lastname, &bloodGroup, &city, &bloodCenter, &regDate)
		if err != nil {
			log.Printf(err.Error())
			return acceptors, err
		}

		acceptors = append(acceptors, Acceptor{
			ID:               id,
			FirstName:        name,
			LastName:         lastname,
			BloodGroup:       bloodGroup,
			City:             city,
			BloodCenter:      bloodCenter,
			RegistrationDate: regDate})
		acceptorsCount++
	}

	return acceptors, nil
}

//GetByID Retrieve an acceptor by Id
func (r *AcceptorsMySQL) GetByID(id string) (Acceptor, error) {
	acceptor := Acceptor{}
	err := r.db.QueryRow(`SELECT * FROM acceptors WHERE id=?`, id).Scan(
		&acceptor.ID,
		&acceptor.FirstName,
		&acceptor.LastName,
		&acceptor.BloodGroup,
		&acceptor.City,
		&acceptor.BloodCenter,
		&acceptor.RegistrationDate)

	return acceptor, err
}

//Update acceptor by ID
func (r *AcceptorsMySQL) Update(acceptor Acceptor) error {
	stmt, err := r.db.Prepare(`UPDATE acceptors SET name=?,lastName=?,city=?,bloodCenter=? WHERE id=?;`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(acceptor.FirstName, acceptor.LastName, acceptor.City, acceptor.BloodCenter, acceptor.ID)
	return err
}

//GetByBloodGroup search for acceptors with specific blood group
func (r *AcceptorsMySQL) GetByBloodGroup(bloodGroup string) ([]Acceptor, error) {
	acceptors := make([]Acceptor, 0)
	rows, err := r.db.Query(`SELECT * FROM acceptors WHERE bloodGroup=?`, bloodGroup)
	if err != nil {
		log.Printf(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, lastname, bloodGroup, city, bloodCenter, regDate string
		err := rows.Scan(&id, &name, &lastname, &bloodGroup, &city, &bloodCenter, &regDate)
		if err != nil {
			log.Printf(err.Error())
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

	return acceptors, err
}

//DeleteByID check whether donor exists and remove
func (r *AcceptorsMySQL) DeleteByID(id string) error {
	_, err := r.db.Query(`DELETE FROM acceptors WHERE id=?`, id)
	return err
}
