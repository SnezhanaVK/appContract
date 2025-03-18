package db

import "github.com/jackc/pgx"

func ConnectDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "1234",
		Database: "contract_db",
		Port:     5432,
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}