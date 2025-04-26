package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)
var (
	dbPool *pgxpool.Pool
)
func SetupDatabase() error {
	connConfig,err := pgx.ParseConfig("host=localhost user=postgres password=1234 dbname=postgres port=5432 sslmode=disable")

	if err != nil {
		return fmt.Errorf("Error parsing database connection string: %v", err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %v", err)
	}
	defer conn.Close(context.Background())

	var exists bool

	err = conn.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'contract_db')`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("Error checking database existence: %v", err)
	}
	if !exists {
		_, err = conn.Exec(context.Background(), `CREATE DATABASE contract_db`)
		if err != nil {
			return fmt.Errorf("Error creating database: %v", err)
		}
		fmt.Println("Database created successfully")
		connConfig.Database= "contract_db"
	conn, err = pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %v", err)
	}
	
	defer conn.Close(context.Background())
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("Error starting transaction: %v", err)
	}
	defer tx.Rollback(context.Background())
	
		_, err = tx.Exec(context.Background(),
	   `CREATE TABLE IF NOT EXISTS roles (
		id_role SERIAL PRIMARY KEY,
		name_role VARCHAR(255) NOT NULL
		)`)
	if err != nil {
	return fmt.Errorf("Error roles creating table: %v", err)
	}
	fmt.Println("Table roles created successfully")
	

	_, err = tx.Exec(context.Background(),`CREATE TABLE IF NOT EXISTS notification_settings (
		id_notification_settings SERIAL PRIMARY KEY,
		variant_notification_settings VARCHAR(30) NOT NULL

	)`)
	if err != nil {
		return fmt.Errorf("Error notification_settings creating table: %v", err)
	}
	fmt.Println("Table notifications created successfully")

	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists users (
		id_user SERIAL PRIMARY KEY,
		surname VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		patronymic VARCHAR(255) NOT NULL,
		phone VARCHAR(20) NOT NULL,
		photo VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		login VARCHAR(255) NOT NULL unique,
		password VARCHAR(255) NOT NULL
		

	)`)
	if err != nil {
		log.Fatal("Error users creating table : ", err)
	}
	fmt.Println("Table users created successfully")
	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists user_by_role (
		user_by_role SERIAL PRIMARY KEY,
		id_user int NOT NULL,
		id_role int NOT NULL,
		CONSTRAINT fk_role_id FOREIGN KEY (id_role) REFERENCES roles(id_role),
		CONSTRAINT fk_user_id FOREIGN KEY (id_user) REFERENCES users(id_user)
	)`)
	if err != nil {
		log.Fatal("Error user_by_role creating table : ", err)
	}
	fmt.Println("Table user_by_role created successfully")

	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists types_contracts (
		id_type_contract SERIAL PRIMARY KEY,
		name_type_contract VARCHAR(255) NOT NULL

	)`)
	if err != nil {
		log.Fatal("Error types_contracts creating table : ", err)
	}
	fmt.Println("Table types_contracts created successfully")

	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists status_contracts (
		id_status_contract SERIAL PRIMARY KEY,
		name_status_contract VARCHAR(255) NOT NULL
		
	)`)
	if err != nil {
		log.Fatal("Error status_contracts creating table : ", err)
	}
	fmt.Println("Table status_contracts created successfully")

	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists counterparty (
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

	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists tegs (
		id_teg SERIAL PRIMARY KEY,
		name_teg VARCHAR(255) NOT NULL
		
	)`)
	if err != nil {
		log.Fatal("Error tegs creating table : ", err)
	}
	fmt.Println("Table tegs created successfully")

	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists status_stages (
		id_status_stage SERIAL PRIMARY KEY,
		name_status_stage VARCHAR(255) NOT NULL
		
	)`)
	if err != nil {
		log.Fatal("Error status_stages creating table : ", err)
	}
	fmt.Println("Table status_stages created successfully")

	_, err = tx.Exec(context.Background(),`CREATE TABLE IF NOT EXISTS contracts (
		id_contract SERIAL PRIMARY KEY,
		name_contract VARCHAR(255) NOT NULL,
		date_create_contract date NOT NULL,
		id_user int NOT NULL,
		date_conclusion date NOT NULL,
		date_end date NOT NULL,
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
	_,err=tx.Exec(context.Background(),`
CREATE TABLE if not exists contract_notifications (
    id_contract_notification SERIAL PRIMARY KEY,
    id_contract int NOT NULL,
    id_user int NOT NULL,
    id_notification_settings int NOT NULL ,
    CONSTRAINT fk_contract FOREIGN KEY (id_contract) REFERENCES contracts(id_contract),
    CONSTRAINT fk_user FOREIGN KEY (id_user) REFERENCES users(id_user),
CONSTRAINT fk_notification_settings FOREIGN KEY (id_notification_settings) REFERENCES notification_settings(id_notification_settings)
);`)

if err != nil {
	log.Fatal("Error contract_notifications creating table : ", err)
}
fmt.Println("Table contract_notifications created successfully")


_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists stages (
    id_stage SERIAL PRIMARY KEY,
    name_stage VARCHAR(255) NOT NULL,
    id_user int NOT NULL,
    description text NOT NULL,
    date_create_start date NOT NULL,
    date_create_end date NOT NULL,
    id_contract int NOT NULL,
    CONSTRAINT id_user FOREIGN KEY (id_user) REFERENCES users(id_user),
    CONSTRAINT id_contract FOREIGN KEY (id_contract) REFERENCES contracts(id_contract) 
)`)
if err != nil {
    log.Fatal("Error stages creating table : ", err)
}
fmt.Println("Table stages created successfully")

_,err=tx.Exec(context.Background(),`CREATE TABLE if not exists stage_notifications (
    id_stage_notification SERIAL PRIMARY KEY,
    id_stage int NOT NULL,
    id_user int NOT NULL,
    id_notification_settings int NOT NULL ,
    CONSTRAINT fk_stage FOREIGN KEY (id_stage) REFERENCES stages(id_stage),
    CONSTRAINT fk_user FOREIGN KEY (id_user) REFERENCES users(id_user),
    CONSTRAINT fk_notification_settings FOREIGN KEY (id_notification_settings) REFERENCES notification_settings(id_notification_settings)
);`)

if err != nil {
	log.Fatal("Error stage_notifications creating table : ", err)
}
fmt.Println("Table stage_notifications created successfully")

_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists history_status (
    id_history_status SERIAL PRIMARY KEY,
    id_stage int NOT NULL,
    id_status_stage int NOT NULL,
    data_change_status date NOT NULL,
    CONSTRAINT id_stage FOREIGN KEY (id_stage) REFERENCES stages(id_stage),
    CONSTRAINT id_status_stage FOREIGN KEY (id_status_stage) REFERENCES status_stages(id_status_stage) 
)`)
if err != nil {
    log.Fatal("Error history_states creating table : ", err)
}
fmt.Println("Table history_states created successfully")

_, err = tx.Exec(context.Background(),`CREATE TABLE IF NOT EXISTS comments (
    id_comment SERIAL PRIMARY KEY,
    id_history_status int NOT NULL,
    comment VARCHAR(1000) NOT NULL,
    date_create_comment date NOT NULL,
 CONSTRAINT id_history_status FOREIGN KEY (id_history_status) REFERENCES history_status(id_history_status)
)`)
if err != nil {
    log.Fatal("Error comments creating table : ", err)
}
fmt.Println("Table comments created successfully")


	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists files (
    id_file SERIAL PRIMARY KEY,
    name_file VARCHAR(255) NOT NULL,
    data bytea NOT NULL,
    type_file VARCHAR(255) NOT NULL,
    id_stage int NOT NULL,
    CONSTRAINT id_stage FOREIGN KEY (id_stage) REFERENCES stages(id_stage) ON DELETE CASCADE
)`)
	
	if err != nil {
		log.Fatal("Error files creating table : ", err)
	}
	fmt.Println("Table files created successfully")

	_, err = tx.Exec(context.Background(),`CREATE TABLE if not exists contracts_by_tegs (
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
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
		
	}

	// Инициализируем пул соединений
	poolConfig, err := pgxpool.ParseConfig("host=localhost user=postgres password=1234 dbname=contract_db port=5432 sslmode=disable")
	if err != nil {
		return fmt.Errorf("error parsing pool config: %v", err)
	}

	dbPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("error creating connection pool: %v", err)
	}

	return nil
}

func GetDBConection() *pgxpool.Pool {
	return dbPool
}

func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
	}
}

