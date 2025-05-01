// File: db/dbConnect.go
package db

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool   *pgxpool.Pool
	initOnce sync.Once
)

func ConnectDB() {
	initOnce.Do(func() {
		config, err := pgxpool.ParseConfig("postgres://postgres:1234@localhost:5432/contract_db")
		if err != nil {
			log.Fatal("Error parsing database config: ", err)
		}

		config.MaxConns = 25
		config.MinConns = 3
		config.MaxConnLifetime = 1 * time.Hour
		config.MaxConnIdleTime = 30 * time.Minute
		config.HealthCheckPeriod = 1 * time.Minute

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Исправлено здесь: используем NewWithConfig вместо ConnectConfig
		dbPool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			log.Fatal("Error connecting to database: ", err)
		}

		log.Println("Successfully connected to database with connection pool")
	})
}

func GetDB() *pgxpool.Pool {
	return dbPool
}

func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
		log.Println("Database connection pool closed")
	}
}