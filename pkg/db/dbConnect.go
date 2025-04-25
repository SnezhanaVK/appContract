package db

// dbConnect.go в папке db
import (
	"log"
	"sync"

	"github.com/jackc/pgx"
)

var (
	pgxConn *pgx.Conn
	initOnce sync.Once

)
func ConnectDB(){
	initOnce.Do(func(){
		var err error
		pgxConn, err = pgx.Connect(pgx.ConnConfig{
			Host:     "localhost",
			User:     "postgres",
			Password: "1234",
			Database: "contract_db",
			Port:     5432,
		})
		if err !=nil{
			log.Fatal("Error conecting to datebase ", err)
		}
		log.Println("Successfully connected to datebase")

	})
}
func GetDB() *pgx.Conn{
	return pgxConn
}
