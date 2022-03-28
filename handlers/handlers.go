package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type User struct {
	Name  string
	Email string
}

type Balance struct {
	ClientID int64
	Balance  int
}

type ApiHandlers struct {
	db    *sql.DB
	cache *redis.Client
}

func NewApiHandlers(db *sql.DB, cache *redis.Client) *ApiHandlers {
	return &ApiHandlers{db, cache}
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

func (h *ApiHandlers) Balances(rw http.ResponseWriter, r *http.Request) {
	log.Println("Query Request")

	res, err := getBalances(h.db, h.cache)
	if err != nil {
		http.Error(rw, "Unable to get balances", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	_, err = rw.Write(res)
	if err != nil {
		http.Error(rw, "Unable to write response", http.StatusInternalServerError)
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

func getBalances(db *sql.DB, cache *redis.Client) ([]byte, error) {
	var data []byte
	key := "balances"
	exp := 30 * time.Second

	value, err := cache.Get(key).Result()
	if err == nil {
		log.Println("Balances loaded from cache")
		return []byte(value), nil
	}

	balances := loadBalances(db)

	data, err = json.Marshal(balances)
	if err != nil {
		return data, err
	}

	err = cache.Set(key, string(data), exp).Err()
	if err != nil {
		return data, err
	}

	log.Println("Balances loaded from db")
	return data, nil
}

func loadBalances(db *sql.DB) []Balance {
	balances := []Balance{}
	var balance Balance

	sql := `
	SELECT client_id, SUM(amount) as balance FROM (
		SELECT * FROM transactions WHERE id NOT IN (
			SELECT DISTINCT transaction_id FROM events WHERE category='cancel'
		)
	) AS txs
	GROUP BY client_id
	ORDER BY balance DESC;
	`
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&balance.ClientID, &balance.Balance)
		if err != nil {
			panic(err)
		}

		balances = append(balances, balance)
	}

	return balances
}
