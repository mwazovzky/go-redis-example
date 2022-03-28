package seeders

import (
	"database/sql"

	"github.com/bxcodec/faker/v3"
)

type Transaction struct {
	Id        int64
	Operation string
	Txid      string `faker:"len=25"`
	Amount    int    `faker:"boundary_start=5, boundary_end=100"`
	ClientId  int64
}

func SeedTransactions(db *sql.DB, clientId int64, count int) ([]Transaction, error) {
	transactions := []Transaction{}

	for i := 0; i < count; i++ {
		tx, err := createTransaction(db, clientId)
		if err != nil {
			panic(err)
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func DeleteTransactions(db *sql.DB) error {
	sql := "DELETE FROM transactions"
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}

	return nil
}

func createTransaction(db *sql.DB, clientId int64) (Transaction, error) {
	operations := []string{"deposit", "withdrawal"}
	tx := Transaction{}

	err := faker.FakeData(&tx)
	if err != nil {
		panic(err)
	}

	tx.ClientId = clientId
	tx.Operation = getRandomArrayItem(operations)
	if tx.Operation == "withdrawal" {
		tx.Amount = -tx.Amount
	}

	sql := "INSERT INTO transactions (operation, txid, amount, client_id) VALUES($1,$2,$3,$4) RETURNING id"
	err = db.QueryRow(sql, tx.Operation, tx.Txid, tx.Amount, tx.ClientId).Scan(&tx.Id)
	if err != nil {
		panic(err)
	}

	return tx, nil
}
