package seeders

import (
	"database/sql"
)

type Event struct {
	Id            int64
	Category      string
	TransactionId int64
}

func SeedEvents(db *sql.DB, transactionId int64, count int) ([]Event, error) {
	events := []Event{}

	for i := 0; i < count; i++ {
		event, err := createEvent(db, transactionId)
		if err != nil {
			panic(err)
		}

		events = append(events, event)
	}

	return events, nil
}

func DeleteEvents(db *sql.DB) error {
	sql := "DELETE FROM events"
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}

	return nil
}

func createEvent(db *sql.DB, transactionId int64) (Event, error) {
	categories := []string{"create", "update", "cancel"}

	event := Event{
		Category:      getRandomArrayItem(categories),
		TransactionId: transactionId,
	}

	sql := "INSERT INTO events (category, transaction_id) VALUES($1,$2)"
	_, err := db.Exec(sql, event.Category, event.TransactionId)
	if err != nil {
		panic(err)
	}

	return event, nil
}
