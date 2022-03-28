package seeders

import (
	"database/sql"

	"github.com/bxcodec/faker/v3"
)

type Client struct {
	Id   int64
	Name string `faker:"sentence"`
}

func SeedClients(db *sql.DB, count int) ([]Client, error) {
	clients := []Client{}

	for i := 0; i < count; i++ {
		client, err := createClient(db)
		if err != nil {
			panic(err)
		}

		clients = append(clients, client)
	}

	return clients, nil
}

func DeleteClients(db *sql.DB) error {
	sql := "DELETE FROM clients"
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}

	return nil
}

func createClient(db *sql.DB) (Client, error) {
	client := Client{}

	err := faker.FakeData(&client)
	if err != nil {
		panic(err)
	}

	sql := "INSERT INTO clients (name) VALUES($1) RETURNING id"
	err = db.QueryRow(sql, client.Name).Scan(&client.Id)
	if err != nil {
		panic(err)
	}

	return client, nil
}
