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
								healthStatus varchar(32),
								donationAmount varchar(32),
								regDate varchar(32),
								lastDonationDate varchar(32),
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

	return db, nil
}
