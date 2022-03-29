package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"go-redis-example/seeders"
)

func init() {
	godotenv.Load()
}

func main() {
	db := connectDB()
	defer db.Close()

	seeders.DeleteClients(db)
	seeders.DeleteTransactions(db)
	seeders.DeleteEvents(db)

	clientCount := 100
	transactionCount := 1000
	eventCount := 5

	clients, _ := seeders.SeedClients(db, clientCount)

	for _, client := range clients {
		transactions, _ := seeders.SeedTransactions(db, client.Id, transactionCount)

		for _, transaction := range transactions {
			seeders.SeedEvents(db, transaction.Id, eventCount)
		}
	}
}

func connectDB() *sql.DB {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected")

	return db
}
