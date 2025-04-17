package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

func SetupDatabase() {
	constStr := "host=localhost user=postgres password=1234 dbname=postgres port=5432 sslmode=disable"

	db, err := sql.Open("postgres", constStr)

	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer db.Close()
// Проверяем, существует ли база данных
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'contract_db')`).Scan(&exists)
	if err != nil {
		log.Fatal("Error checking database existence: ", err)
	}
	if !exists {
	_, err = db.Exec(`CREATE DATABASE contract_db`)
	if err != nil {
		log.Fatal("Error creating database: ", err)
	}
	fmt.Println("Database created successfully")
	connStr := "host=localhost user=postgres password=1234 dbname=contract_db port=5432 sslmode=disable"
	
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS roles (
		id_role SERIAL PRIMARY KEY,
		name_role VARCHAR(255) NOT NULL
		)`)
	if err != nil {
		log.Fatal("Error roles creating table : ", err)
	}
	fmt.Println("Table roles created successfully")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS notifications (
		id_notification SERIAL PRIMARY KEY,
		variant_notification VARCHAR(30) NOT NULL

	)`)
	if err != nil {
		log.Fatal("Error notifications creating table : ", err)
	}
	fmt.Println("Table notifications created successfully")

	_, err = db.Exec(`CREATE TABLE if not exists users (
		id_user SERIAL PRIMARY KEY,
		surname VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		patronymic VARCHAR(255) NOT NULL,
		phone VARCHAR(20) NOT NULL,
		photo VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		role_id int NOT NULL,
		notification_id int NOT NULL,
		admin bool NOT NULL,
		login VARCHAR(255) NOT NULL unique,
		password VARCHAR(255) NOT NULL,
		CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES roles(id_role),
		CONSTRAINT fk_notification_id FOREIGN KEY (notification_id) REFERENCES notifications(id_notification)

	)`)
	if err != nil {
		log.Fatal("Error users creating table : ", err)
	}
	fmt.Println("Table users created successfully")

	_, err = db.Exec(`CREATE TABLE if not exists types_contracts (
		id_type_contract SERIAL PRIMARY KEY,
		name_type_contract VARCHAR(255) NOT NULL

	)`)
	if err != nil {
		log.Fatal("Error types_contracts creating table : ", err)
	}
	fmt.Println("Table types_contracts created successfully")

	_, err = db.Exec(`CREATE TABLE if not exists status_contracts (
		id_status_contract SERIAL PRIMARY KEY,
		name_status_contract VARCHAR(255) NOT NULL
		
	)`)
	if err != nil {
		log.Fatal("Error status_contracts creating table : ", err)
	}
	fmt.Println("Table status_contracts created successfully")

	_, err = db.Exec(`CREATE TABLE if not exists counterparty (
		id_counterparty SERIAL PRIMARY KEY,
		name_counterparty VARCHAR(255) NOT NULL,
		contact VARCHAR(255) NOT NULL,
		inn VARCHAR(255) NOT NULL,
		ogrn VARCHAR(255) NOT NULL,
		address VARCHAR(255) NOT NULL,
		dop_info VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		log.Fatal("Error counterparty creating table : ", err)
	}
	fmt.Println("Table counterparty created successfully")

	_, err = db.Exec(`CREATE TABLE if not exists tegs (
		id_teg SERIAL PRIMARY KEY,
		name_teg VARCHAR(255) NOT NULL
		
	)`)
	if err != nil {
		log.Fatal("Error tegs creating table : ", err)
	}
	fmt.Println("Table tegs created successfully")

	_, err = db.Exec(`CREATE TABLE if not exists status_stages (
		id_status_stage SERIAL PRIMARY KEY,
		name_status_stage VARCHAR(255) NOT NULL
		
	)`)
	if err != nil {
		log.Fatal("Error status_stages creating table : ", err)
	}
	fmt.Println("Table status_stages created successfully")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS contracts (
		id_contract SERIAL PRIMARY KEY,
		name_contract VARCHAR(255) NOT NULL,
		date_create_contract date NOT NULL,
		id_user int NOT NULL,
		data_conclusion date NOT NULL,
		data_end date NOT NULL,
		id_type int NOT NULL,
		cost int NOT NULL,
		object_contract VARCHAR(255) NOT NULL,
		term_payment varchar(255) NOT NULL, 
		id_counterparty int NOT NULL,
		id_status_contract int NOT NULL,
		notes varchar(1000) NOT NULL,
		conditions varchar(1000) NOT NULL,
		CONSTRAINT id_user FOREIGN KEY (id_user) REFERENCES users(id_user),
		CONSTRAINT id_type FOREIGN KEY (id_type) REFERENCES types_contracts(id_type_contract),
		CONSTRAINT id_counterparty FOREIGN KEY (id_counterparty) REFERENCES counterparty(id_counterparty),
		CONSTRAINT id_status_contract FOREIGN KEY (id_status_contract) REFERENCES status_contracts(id_status_contract)
	)`)
	if err != nil {
		log.Fatal("Error contracts creating table : ", err)
	}
	fmt.Println("Table contracts created successfully")
	

	
_, err = db.Exec(`CREATE TABLE if not exists stages (
    id_stage SERIAL PRIMARY KEY,
    name_stage VARCHAR(255) NOT NULL,
    id_user int NOT NULL,
    description text NOT NULL,
    id_status_stage int NOT NULL,
    date_create_start date NOT NULL,
    date_create_end date NOT NULL,
    id_contract int NOT NULL,
    CONSTRAINT id_user FOREIGN KEY (id_user) REFERENCES users(id_user),
    CONSTRAINT id_status_stage FOREIGN KEY (id_status_stage) REFERENCES status_stages(id_status_stage),
    CONSTRAINT id_contract FOREIGN KEY (id_contract) REFERENCES contracts(id_contract)
)`)
	if err != nil {
		log.Fatal("Error stages creating table : ", err)
	}
	fmt.Println("Table stages created  successfully")

	_, err = db.Exec(`CREATE TABLE if not exists files (
    id_file SERIAL PRIMARY KEY,
    name_file VARCHAR(255) NOT NULL,
    data byte NOT NULL,
    type_file VARCHAR(255) NOT NULL,
    id_stage int NOT NULL,
    CONSTRAINT id_stage FOREIGN KEY (id_stage) REFERENCES stages(id_stage) ON DELETE CASCADE
)`)
	
	if err != nil {
		log.Fatal("Error files creating table : ", err)
	}
	fmt.Println("Table files created successfully")

	_, err = db.Exec(`CREATE TABLE if not exists history_states (
		id_history_state SERIAL PRIMARY KEY,
		id_stage int NOT NULL,
		id_status_stage int NOT NULL,
		data_create date NOT NULL,
		comment text NOT NULL,
		CONSTRAINT id_stage FOREIGN KEY (id_stage) REFERENCES stages(id_stage),
		CONSTRAINT id_status_stage FOREIGN KEY (id_status_stage) REFERENCES status_stages(id_status_stage)
	)`)
	if err != nil {
		log.Fatal("Error history_states creating table : ", err)
	}
	fmt.Println("Table history_states created successfully")

	_, err = db.Exec(`CREATE TABLE if not exists contracts_by_tegs (
		id_contract_by_teg SERIAL PRIMARY KEY,
		id_contract int NOT NULL,
		
		id_teg int NOT NULL,
		CONSTRAINT id_contract FOREIGN KEY (id_contract) REFERENCES contracts(id_contract),
		CONSTRAINT id_teg FOREIGN KEY (id_teg) REFERENCES tegs(id_teg)
	)`)
	if err != nil {
		log.Fatal("Error contracts_by_tegs creating table : ", err)
	}
	fmt.Println("Table contracts_by_tegs created successfully")
}

}
