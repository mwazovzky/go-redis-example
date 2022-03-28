package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Name  string
	Email string
}

type ApiHandlers struct {
	db *sql.DB
}

func NewApiHandlers(db *sql.DB) *ApiHandlers {
	return &ApiHandlers{db}
}

func (h *ApiHandlers) Users(rw http.ResponseWriter, r *http.Request) {
	log.Println("User Index Request")

	users := getUsers(h.db)

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err := e.Encode(users)

	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func getUsers(db *sql.DB) []User {
	users := []User{}
	var user User

	rows, err := db.Query("SELECT name, email FROM users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&user.Name, &user.Email)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	return users
}
