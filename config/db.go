package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func CreateDatabaseConn() (*sql.DB, error) {
	serverName := "localhost:3306"
	user := "docker"
	password := "password"
	dbName := "accounts"
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, serverName, dbName)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("USE accounts")
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("DB selected successfully...")
	}

	stmt, err := db.Prepare(`DROP TABLE IF EXISTS donors, acceptors`)

	if err != nil {
		log.Printf(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Acceptors table dropped successfully...")
		log.Printf("Donors table dropped successfully...")
	}

	stmtAcceptors, err := db.Prepare(`CREATE TABLE acceptors (
								id int NOT NULL,
								name varchar(32),
								lastName varchar(32),
								bloodGroup varchar(32),
								city varchar(50),
								bloodCenter varchar(100),
								regDate varchar(32),
								healthStatus varchar(32),
								PRIMARY KEY (id));`)
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = stmtAcceptors.Exec()
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Acceptors table created successfully...")
	}

	stmtDonors, err := db.Prepare(`CREATE TABLE donors (
								id int NOT NULL,
								name varchar(32),
								lastName varchar(32),
								phone varchar(32),
								email varchar(32),
								age integer,
								gender varchar(32),
								bloodGroup varchar(32),
								city varchar(50),
								bloodCenter varchar(100),
								regDate varchar(32),
								PRIMARY KEY (id)
						);`)

	if err != nil {
		log.Printf(err.Error())
	}
	_, err = stmtDonors.Exec()
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Donors table created successfully...")
	}

	mockdata, err := db.Prepare(`INSERT INTO donors(id, name, lastName, phone, email, age, gender, bloodGroup, city, bloodCenter, regDate)
								VALUES ('12','Ivan','Petrov','08978654321','ivanp@abv.bg','31','MALE','AB+','Sofia', 'Pirogof', 'Sun Mar 15 02:44:15 EET 2019');`)
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = mockdata.Exec()
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Mock data added...")
	}

	return db, nil
}
