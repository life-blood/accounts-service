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
								id varchar(32) NOT NULL,
								name varchar(32),
								lastName varchar(32),
								bloodGroup varchar(32),
								city varchar(50),
								bloodCenter varchar(250),
								regDate varchar(32),
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
								id varchar(32) NOT NULL,
								name varchar(32),
								lastName varchar(32),
								phone varchar(32),
								email varchar(32),
								age integer,
								gender varchar(32),
								bloodGroup varchar(32),
								city varchar(50),
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

	mockdata1, err := db.Prepare(`INSERT INTO donors(id, name, lastName, phone, email, age, gender, bloodGroup, city, regDate)
								VALUES ('12','Ivan','Petrov','08978654321','ivanp@abv.bg','31','MALE','AB','Sofia', 'Sun Mar 15 02:44:15 EET 2019');`)
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = mockdata1.Exec()
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Mock data 1 added...")
	}

	mockdata2, err := db.Prepare(`INSERT INTO acceptors(id, name, lastName, bloodGroup, city, bloodCenter, regDate)
								VALUES ('12','Ivan','Petrov','AB','Sofia', 'РЦ по трансфузионна хематология - Пловдив', 'Sun Mar 15 02:44:15 EET 2019');`)
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = mockdata2.Exec()
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Mock data 2 added...")
	}

	mockdata3, err := db.Prepare(`INSERT INTO donors(id, name, lastName, phone, email, age, gender, bloodGroup, city, regDate)
								VALUES ('1','Petka','Petrova','08978654321','ivanp@abv.bg','31','MALE','B','Sofia', 'Sun Mar 15 02:44:15 EET 2019');`)
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = mockdata3.Exec()
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Mock data 3 added...")
	}

	mockdata4, err := db.Prepare(`INSERT INTO acceptors(id, name, lastName, bloodGroup, city, bloodCenter, regDate)
								VALUES ('2','Ivaylo','Yosifov','0','Plovdiv', 'РЦ по трансфузионна хематология - Варна', 'Sun Mar 16 02:44:15 EET 2020');`)
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = mockdata4.Exec()
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Mock data 4 added...")
	}

	return db, nil
}
