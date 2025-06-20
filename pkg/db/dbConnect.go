package db

import (
	"context"
	"log"
	"os"
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
		
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "postgres")
		dbPass := getEnv("DB_PASSWORD", "1234")
		dbName := getEnv("DB_NAME", "contract_db")
		sslMode := getEnv("SSL_MODE", "disable") 

		
		dsn := "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + sslMode

		config, err := pgxpool.ParseConfig(dsn)
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
		
		dbPool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			log.Fatal("Error connecting to database: ", err)
		}

		log.Println("Successfully connected to database with connection pool")
	})
}


func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
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
