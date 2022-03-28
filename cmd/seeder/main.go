package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"go-redis-example/seeders"
)

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
	dbhost := "postgres"
	dbport := 5432
	dbuser := "user"
	dbpassword := "secret"
	dbname := "testdb"

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", conn)
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
